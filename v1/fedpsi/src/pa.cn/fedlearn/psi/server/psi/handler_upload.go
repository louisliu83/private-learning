package psi

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"pa.cn/fedlearn/psi/api"
	"pa.cn/fedlearn/psi/api/types"
	log "pa.cn/fedlearn/psi/log"
)

func (s *Server) CheckChunk(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "CheckChunk")
	log.Debugln(ctx, "Server.CheckChunk is called.")

	var r types.ChunkCheckRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "CheckChunk", "Failed")
		api.OutputError(res, err)
		return
	}

	if ok, err := s.UploadMgr.CheckChunk(ctx, r); err != nil {
		auditContextActionResult(ctx, "CheckChunk", "Failed")
		api.OutputError(res, err)
		return
	} else {
		if !ok {
			auditContextActionResult(ctx, "CheckChunk", "Failed")
			api.OutputError(res, errors.New("False"))
			return
		}
	}

	auditContextActionResult(ctx, "CheckChunk", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) UploadChunk(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "UploadChunk")
	log.Debugln(ctx, "Server.UploadChunk is called.")

	var r types.ChunkUploadRequest
	maxUploadSize := int64(64 * 1024 * 1024)
	req.Body = http.MaxBytesReader(res, req.Body, maxUploadSize)

	err := req.ParseMultipartForm(maxUploadSize)
	if err != nil {
		auditContextActionResult(ctx, "UploadChunk", "Failed")
		api.OutputError(res, err)
		return
	}

	file, _, err := req.FormFile("file")
	if err != nil {
		auditContextActionResult(ctx, "UploadChunk", "Failed")
		api.OutputError(res, err)
		return
	}

	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		auditContextActionResult(ctx, "UploadChunk", "Failed")
		api.OutputError(res, err)
		return
	}

	r.FileData = fileBytes
	r.MD5 = req.FormValue("md5")
	strChunk := req.FormValue("chunk")
	if strChunk == "" || strChunk == "null" {
		strChunk = "0"
	}

	if chunk, err := strconv.Atoi(strChunk); err != nil {
		log.Errorln(ctx, "Invalid chunk index")
		auditContextActionResult(ctx, "UploadChunk", "Failed")
		api.OutputError(res, err)
		return
	} else {
		r.Chunk = int32(chunk)
	}

	err = s.UploadMgr.UploadChunk(ctx, r)
	if err != nil {
		auditContextActionResult(ctx, "UploadChunk", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "UploadChunk", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) MergeChunk(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "MergeChunk")
	log.Debugln(ctx, "Server.MergeChunk is called.")

	var r types.ChunkMergeRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "MergeChunk", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.UploadMgr.MergeChunk(ctx, r); err != nil {
		auditContextActionResult(ctx, "MergeChunk", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "MergeChunk", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) FilePull(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "FilePull")
	log.Debugln(ctx, "Server.FilePull is called.")

	var r types.FilePullRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "FilePull", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.UploadMgr.FilePull(ctx, r); err != nil {
		auditContextActionResult(ctx, "FilePull", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "FilePull", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) FileCopy(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "FileCopy")
	log.Debugln(ctx, "Server.FileCopy is called.")

	var r types.FileLocalCopyRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "FileCopy", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.UploadMgr.FileLocalCopy(ctx, r); err != nil {
		auditContextActionResult(ctx, "FileCopy", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "FileCopy", "Success")
	api.OutputSuccess(res)
	return
}
