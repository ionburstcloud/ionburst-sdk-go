package ionburst

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/google/uuid"
	"gitlab.com/ionburst/ionburst-sdk-go/models"
)

type DeferredToken string

const objectLimit int = 50000000

func (cli *Client) GetClassifications() ([]string, error) {
	cli.logger.Debug("Getting Classifications")
	if res, err := cli.doGet("api/classification", nil); err != nil {
		return nil, err
	} else {
		var classifications []string
		if err := json.Unmarshal(res.Body(), &classifications); err != nil {
			return nil, err
		}
		cli.logger.Debug("Retrieved Classifications")
		return classifications, nil
	}
}

func (cli *Client) GetDataUploadLimit() (int, error) {
	cli.logger.Debug("Getting Data Upload Limit")
	if res, err := cli.doGet("api/data/query/uploadsizelimit", nil); err != nil {
		return 0, err
	} else {
		var datalimitbody []byte = res.Body()
		var datauploadlimit, _ = strconv.Atoi(string(datalimitbody))

		cli.logger.Debug("Retrieved Data Upload Limit")
		return datauploadlimit, nil
	}
}

func (cli *Client) GetSecretsUploadLimit() (int, error) {
	cli.logger.Debug("Getting Secrets Upload Limit")
	if res, err := cli.doGet("api/secrets/query/uploadsizelimit", nil); err != nil {
		return 0, err
	} else {
		var secretslimitbody []byte = res.Body()
		var secretsuploadlimit, _ = strconv.Atoi(string(secretslimitbody))

		cli.logger.Debug("Retrieved Secrets Upload Limit")
		return secretsuploadlimit, nil
	}
}

func (cli *Client) Put(id string, reader io.Reader, classification string) error {
	cli.logger.Debug("Uploading Ionburst object: ", id)
	_, err := cli.doPostBinary("api/data/"+id, reader, classification)
	if err == nil {
		cli.logger.Debug("Ionburst Upload Complete for: ", id)
	}
	return err
}

func (cli *Client) PutSecrets(id string, reader io.Reader, classification string) error {
	cli.logger.Debug("Uploading Ionburst secret: ", id)
	_, err := cli.doPostBinary("api/secrets/"+id, reader, classification)
	if err == nil {
		cli.logger.Debug("Ionburst Put Complete for: ", id)
	}
	return err
}

func (cli *Client) PutFromFile(id string, file string, classification string) error {
	cli.logger.Debug("Uploading Ionburst object ", id, " from file ", file)
	rdr, err := os.Open(file)
	if err != nil {
		return err
	}
	_, err = cli.doPostBinary("api/data/"+id, rdr, classification)
	if err == nil {
		cli.logger.Debug("Ionburst Upload Complete for: ", id)
	}
	return err
}

func (cli *Client) PutSecretsFromFile(id string, file string, classification string) error {
	cli.logger.Debug("Uploading Ionburst secret ", id, " from file ", file)
	rdr, err := os.Open(file)
	if err != nil {
		return err
	}
	_, err = cli.doPostBinary("api/secrets/"+id, rdr, classification)
	if err == nil {
		cli.logger.Debug("Ionburst Upload Complete for: ", id)
	}
	return err
}

func (cli *Client) PutManifest(id string, reader io.Reader, classification string) error {
	cli.logger.Debug("Creating manifest for Ionburst object: ", id)

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return err
	}

	hash, err := objectHash(reader)

	fileInfo := buf.Len()

	var fileSize = int64(fileInfo)

	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(objectLimit)))

	cli.logger.Debug("Number of chunks: ", totalPartsNum)

	var manifest models.Manifest

	manifest.Name = id
	manifest.ChunkCount = int(totalPartsNum)
	manifest.ChunkSize = objectLimit
	manifest.MaxSize = objectLimit
	manifest.Size = int(fileSize)
	manifest.Hash = hash

	for i := uint64(0); i < totalPartsNum; i++ {
		chunkID := uuid.NewString()
		cli.logger.Debug("Chunk ID: ", chunkID)

		partSize := int64(math.Min(float64(objectLimit), float64(fileSize-int64(i*uint64(objectLimit)))))
		cli.logger.Debug("Chunk size: ", partSize)
		offset := int64(i * uint64(objectLimit))
		cli.logger.Debug("Chunk offset: ", offset)
		buffer := make([]byte, partSize)

		_, err := buf.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			fmt.Println(err)
		}

		r := bytes.NewReader(buffer)
		chunkHash, err := objectHash(r)

		r = bytes.NewReader(buffer)
		_, err = cli.doPostBinary("api/data/"+chunkID, r, classification)
		if err != nil {
			cli.logger.Debug(err)
		} else if err == nil {
			cli.logger.Debug("Ionburst Chunk Upload Complete for: ", chunkID)
		}

		fmt.Println("Split to: ", chunkID)

		manifest.Chunks = append(manifest.Chunks, models.Chunks{
			ID:   chunkID,
			Ord:  int(i) + 1,
			Hash: chunkHash,
		})
	}
	m, err := json.Marshal(manifest)

	r := bytes.NewReader(m)

	_, err = cli.doPostBinary("api/data/"+id, r, classification)
	if err == nil {
		cli.logger.Debug("Ionburst Manifest Upload Complete for: ", id)
	}

	return err

}

func (cli *Client) PutManifestFromFile(id string, file string, classification string) error {
	cli.logger.Debug("Creating manifest for Ionburst object: ", id, " from file ", file)
	rdr, err := os.Open(file)
	if err != nil {
		return err
	}

	hash, err := objectHash(rdr)

	fileInfo, _ := rdr.Stat()

	var fileSize = fileInfo.Size()

	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(objectLimit)))

	cli.logger.Debug("Number of chunks: ", totalPartsNum)

	var manifest models.Manifest

	manifest.Name = id
	manifest.ChunkCount = int(totalPartsNum)
	manifest.ChunkSize = objectLimit
	manifest.MaxSize = objectLimit
	manifest.Size = int(fileSize)
	manifest.Hash = hash

	for i := uint64(0); i < totalPartsNum; i++ {
		chunkID := uuid.NewString()
		cli.logger.Debug("Chunk ID: ", chunkID)

		partSize := int64(math.Min(float64(objectLimit), float64(fileSize-int64(i*uint64(objectLimit)))))
		cli.logger.Debug("Chunk size: ", partSize)
		offset := int64(i * uint64(objectLimit))
		cli.logger.Debug("Chunk offset: ", offset)
		buffer := make([]byte, partSize)

		_, err := rdr.ReadAt(buffer, offset)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			fmt.Println("whoops")
			fmt.Println(err)
			break
		}

		r := bytes.NewReader(buffer)
		chunkHashRaw := sha256.New()

		if _, err := io.Copy(chunkHashRaw, r); err != nil {
			log.Println(err)
		}

		chunkHash := base64.StdEncoding.EncodeToString(chunkHashRaw.Sum(nil))

		r = bytes.NewReader(buffer)
		_, err = cli.doPostBinary("api/data/"+chunkID, r, classification)
		if err != nil {
			cli.logger.Debug(err)
		} else if err == nil {
			cli.logger.Debug("Ionburst Chunk Upload Complete for: ", chunkID)
		}

		fmt.Println("Split to: ", chunkID)

		manifest.Chunks = append(manifest.Chunks, models.Chunks{
			ID:   chunkID,
			Ord:  int(i) + 1,
			Hash: chunkHash,
		})
	}
	m, err := json.Marshal(manifest)

	r := bytes.NewReader(m)

	_, err = cli.doPostBinary("api/data/"+id, r, classification)
	if err == nil {
		cli.logger.Debug("Ionburst Manifest Upload Complete for: ", id)
	}

	return err
}

func (cli *Client) Get(id string) (io.Reader, error) {
	cli.logger.Debug("Downloading Ionburst object: ", id)
	return cli.doGetBinary("api/data/"+id, nil)
}

func (cli *Client) GetSecrets(id string) (io.Reader, error) {
	cli.logger.Debug("Downloading Ionburst secret: ", id)
	return cli.doGetBinary("api/secrets/"+id, nil)
}

func (cli *Client) GetWithLen(id string) (io.Reader, int64, error) {
	cli.logger.Debug("Downloading Ionburst object: ", id)
	return cli.doGetBinaryLen("api/data/"+id, nil)
}

func (cli *Client) GetSecretsWithLen(id string) (io.Reader, int64, error) {
	cli.logger.Debug("Downloading Ionburst secret: ", id)
	return cli.doGetBinaryLen("api/secrets/"+id, nil)
}

func (cli *Client) GetToFile(id string, file string) error {
	cli.logger.Debug("Downloading Ionburst object: ", id, " to file ", file)
	rdr, err := cli.doGetBinary("api/data/"+id, nil)
	if err != nil {
		return err
	}
	wr, err := os.Create(file)
	if err != nil {
		return err
	}
	_, err = io.Copy(wr, rdr)
	cli.logger.Debug("Ionburst object download: ", id, " complete")
	return err
}

func (cli *Client) GetSecretsToFile(id string, file string) error {
	cli.logger.Debug("Downloading Ionburst secret: ", id, " to file ", file)
	rdr, err := cli.doGetBinary("api/secrets/"+id, nil)
	if err != nil {
		return err
	}
	wr, err := os.Create(file)
	if err != nil {
		return err
	}
	_, err = io.Copy(wr, rdr)
	cli.logger.Debug("Ionburst secret download: ", id, " complete")
	return err
}

func (cli *Client) GetManifest(id string) (io.Reader, error) {
	cli.logger.Debug("Starting download of Ionburst manifest: ", id)

	cli.logger.Debug("Retrieving Ionburst manifest: ", id)

	manifestObject, err := cli.doGetBinary("api/data/"+id, nil)
	if err != nil {
		return nil, err
	}

	var manifest models.Manifest

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(manifestObject)
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(buf.Bytes(), &manifest)

	objectBuf := new(bytes.Buffer)

	for i := 0; i < len(manifest.Chunks); i++ {
		cli.logger.Debug("Retrieving chunk: ", manifest.Chunks[i].ID)
		chunk, err := cli.doGetBinary("api/data/"+manifest.Chunks[i].ID, nil)
		if err != nil {
			return nil, err
		} else {
			_, err := objectBuf.ReadFrom(chunk)
			if err != nil {
				return nil, err
			}
		}
		cli.logger.Debug("Retrieval of chunk: ", manifest.Chunks[i].ID, " complete")
	}

	reader := bytes.NewReader(objectBuf.Bytes())

	return reader, err
}

func (cli *Client) GetManifestToFile(id string, file string) error {
	cli.logger.Debug("Starting download of Ionburst manifest: ", id)

	cli.logger.Debug("Retrieving Ionburst manifest: ", id)

	manifestObject, err := cli.doGetBinary("api/data/"+id, nil)
	if err != nil {
		return err
	}

	var manifest models.Manifest

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(manifestObject)
	if err != nil {
		return err
	}

	_ = json.Unmarshal(buf.Bytes(), &manifest)

	objectBuf := new(bytes.Buffer)

	for i := 0; i < len(manifest.Chunks); i++ {
		cli.logger.Debug("Retrieving chunk: ", manifest.Chunks[i].ID)
		chunk, err := cli.doGetBinary("api/data/"+manifest.Chunks[i].ID, nil)
		if err != nil {
			return err
		} else {
			_, err := objectBuf.ReadFrom(chunk)
			if err != nil {
				return err
			}
		}
		cli.logger.Debug("Retrieval of chunk: ", manifest.Chunks[i].ID, " complete")
	}

	reader := bytes.NewReader(objectBuf.Bytes())

	wr, err := os.Create(file)
	if err != nil {
		return err
	}
	_, err = io.Copy(wr, reader)
	cli.logger.Debug("Ionburst manifest download ", id, " complete")

	return err
}

func (cli *Client) Delete(id string) error {
	cli.logger.Debug("Starting Deletion of Ionburst object: ", id)
	_, err := cli.doDelete("api/data/"+id, nil)
	cli.logger.Debug("Deletion of Ionburst object: ", id, " complete")
	return err
}

func (cli *Client) DeleteSecrets(id string) error {
	cli.logger.Debug("Starting Deletion of Ionburst secret: ", id)
	_, err := cli.doDelete("api/secrets/"+id, nil)
	cli.logger.Debug("Deletion of Ionburst secret: ", id, " complete")
	return err
}

func (cli *Client) DeleteManifest(id string) error {
	cli.logger.Debug("Starting Deletion of Ionburst manifest: ", id)

	cli.logger.Debug("Retrieving Ionburst manifest: ", id)

	manifestObject, err := cli.doGetBinary("api/data/"+id, nil)
	if err != nil {
		return err
	}

	var manifest models.Manifest

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(manifestObject)
	if err != nil {
		return err
	}

	_ = json.Unmarshal(buf.Bytes(), &manifest)

	for i := 0; i < len(manifest.Chunks); i++ {
		cli.logger.Debug("Deleting chunk: ", manifest.Chunks[i].ID)
		_, err := cli.doDelete("api/data/"+manifest.Chunks[i].ID, nil)
		if err != nil {
			return err
		}
		cli.logger.Debug("Deletion of chunk: ", manifest.Chunks[i].ID, " complete")
	}

	cli.logger.Debug("Deleting manifest: ", id)
	_, err = cli.doDelete("api/data/"+id, nil)

	return err
}

func (cli *Client) Head(id string) error {
	cli.logger.Debug("Querying Ionburst object: ", id)
	_, err := cli.doHead("api/data/"+id, nil)
	cli.logger.Debug("Query of Ionburst object: ", id, " complete")
	return err
}

func (cli *Client) HeadSecrets(id string) error {
	cli.logger.Debug("Querying Ionburst secret: ", id)
	_, err := cli.doHead("api/secrets/"+id, nil)
	cli.logger.Debug("Query of Ionburst secret: ", id, " complete")
	return err
}

func (cli *Client) HeadWithLen(id string) (int64, error) {
	cli.logger.Debug("Querying Ionburst object: ", id)
	return cli.doHeadLen("api/data/"+id, nil)
}

func (cli *Client) HeadSecretsWithLen(id string) (int64, error) {
	cli.logger.Debug("Querying Ionburst secret: ", id)
	return cli.doHeadLen("api/secrets/"+id, nil)
}

func (cli *Client) PutDeferred(id string, reader io.Reader, classification string) (DeferredToken, error) {
	cli.logger.Debug("Uploading Ionburst object deferred for: ", id)
	tk, err := cli.doPostBinaryDeferred("api/data/deferred/start/"+id, reader, classification)
	if err == nil {
		cli.logger.Debug("Ionburst Deferred Upload Started for: ", id, " token: ", tk)
	} else {
		return "", err
	}
	return DeferredToken(tk), nil
}

func (cli *Client) PutDeferredSecrets(id string, reader io.Reader, classification string) (DeferredToken, error) {
	cli.logger.Debug("Uploading Ionburst secret deferred for: ", id)
	tk, err := cli.doPostBinaryDeferred("api/secrets/deferred/start/"+id, reader, classification)
	if err == nil {
		cli.logger.Debug("Ionburst Deferred Upload Started for: ", id, " token: ", tk)
	} else {
		return "", err
	}
	return DeferredToken(tk), nil
}

func (cli *Client) GetDeferred(id string) (DeferredToken, error) {
	cli.logger.Debug("Downloading Ionburst object deferred for: ", id)
	tk, err := cli.doGet("api/data/deferred/start/"+id, nil)
	if err == nil {

		var defResp models.DeferredTokenResponse
		err = json.Unmarshal(tk.Body(), &defResp)
		if err != nil {
			return "", err
		}
		cli.logger.Debug("Ionburst Deferred Upload Started for: ", id, " token: ", defResp.DeferredToken)
		return DeferredToken(defResp.DeferredToken), nil
	} else {
		return "", err
	}
}

func (cli *Client) GetDeferredSecrets(id string) (DeferredToken, error) {
	cli.logger.Debug("Downloading Ionburst secret deferred for: ", id)
	tk, err := cli.doGet("api/secrets/deferred/start/"+id, nil)
	if err == nil {

		var defResp models.DeferredTokenResponse
		err = json.Unmarshal(tk.Body(), &defResp)
		if err != nil {
			return "", err
		}
		cli.logger.Debug("Ionburst Deferred GET Started for: ", id, " token: ", defResp.DeferredToken)
		return DeferredToken(defResp.DeferredToken), nil
	} else {
		return "", err
	}
}

func (cli *Client) CheckDeferred(token DeferredToken) (*models.WorkflowResult, error) {
	res, err := cli.doGet("api/data/deferred/check/"+string(token), nil)
	if err != nil {
		return nil, err
	}
	var wr models.WorkflowResult
	err = json.Unmarshal(res.Body(), &wr)
	if err != nil {
		return nil, err
	}
	return &wr, nil
}

func (cli *Client) CheckDeferredSecrets(token DeferredToken) (*models.WorkflowResult, error) {
	res, err := cli.doGet("api/secrets/deferred/check/"+string(token), nil)
	if err != nil {
		return nil, err
	}
	var wr models.WorkflowResult
	err = json.Unmarshal(res.Body(), &wr)
	if err != nil {
		return nil, err
	}
	return &wr, nil
}

func (cli *Client) FetchDeferred(token DeferredToken) (io.Reader, error) {
	cli.logger.Debug("Downloading Deferred Ionburst object with token: ", token)
	return cli.doGetBinary("api/data/deferred/fetch/"+string(token), nil)
}

func (cli *Client) FetchDeferredSecrets(token DeferredToken) (io.Reader, error) {
	cli.logger.Debug("Downloading Deferred Ionburst secret with token: ", token)
	return cli.doGetBinary("api/secrets/deferred/fetch/"+string(token), nil)
}
