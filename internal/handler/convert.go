package handler

import (
	"archive/zip"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"kgm2flac-backend/internal/config"
	"kgm2flac-backend/internal/service"
	"kgm2flac-backend/internal/utils"
	"kgm2flac-backend/pkg/types"
)

const page = `<!doctype html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>KGM → FLAC 转换器</title>
    <style>
        /* CSS 样式保持不变 */
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
            color: #333;
        }
        
        .container {
            max-width: 800px;
            margin: 0 auto;
            background: white;
            border-radius: 12px;
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
            overflow: hidden;
        }
        
        .header {
            background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
            color: white;
            padding: 30px;
            text-align: center;
        }
        
        .header h1 {
            font-size: 2.2em;
            margin-bottom: 10px;
            font-weight: 300;
        }
        
        .header p {
            opacity: 0.9;
            font-size: 1.1em;
        }
        
        .content {
            padding: 40px;
        }
        
        .upload-area {
            border: 3px dashed #4facfe;
            border-radius: 8px;
            padding: 40px;
            text-align: center;
            background: #f8f9fa;
            transition: all 0.3s ease;
            margin-bottom: 30px;
        }
        
        .upload-area:hover {
            border-color: #00f2fe;
            background: #e3f2fd;
        }
        
        .upload-area.dragover {
            border-color: #00c853;
            background: #e8f5e8;
        }
        
        .upload-icon {
            font-size: 3em;
            color: #4facfe;
            margin-bottom: 15px;
        }
        
        .file-input {
            display: none;
        }
        
        .browse-btn {
            background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
            color: white;
            padding: 12px 24px;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-size: 1em;
            transition: transform 0.2s ease;
            margin: 10px 0;
        }
        
        .browse-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(79, 172, 254, 0.3);
        }
        
        .submit-btn {
            background: linear-gradient(135deg, #00b09b 0%, #96c93d 100%);
            color: white;
            padding: 15px 30px;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-size: 1.1em;
            font-weight: 500;
            transition: all 0.3s ease;
            width: 100%;
            margin-top: 20px;
        }
        
        .submit-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 20px rgba(0, 176, 155, 0.4);
        }
        
        .submit-btn:disabled {
            background: #ccc;
            cursor: not-allowed;
            transform: none;
            box-shadow: none;
        }
        
        .file-list {
            margin-top: 20px;
            max-height: 200px;
            overflow-y: auto;
        }
        
        .file-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 10px;
            background: #f8f9fa;
            border-radius: 6px;
            margin-bottom: 8px;
            border-left: 4px solid #4facfe;
        }
        
        .file-name {
            flex: 1;
            font-weight: 500;
        }
        
        .file-size {
            color: #666;
            font-size: 0.9em;
            margin-left: 10px;
        }
        
        .remove-btn {
            background: #ff4757;
            color: white;
            border: none;
            border-radius: 4px;
            padding: 4px 8px;
            cursor: pointer;
            font-size: 0.8em;
        }
        
        .remove-btn:hover {
            background: #ff3742;
        }
        
        .progress-container {
            margin-top: 20px;
        }
        
        .progress-bar {
            width: 100%;
            height: 6px;
            background: #e0e0e0;
            border-radius: 3px;
            overflow: hidden;
        }
        
        .progress-fill {
            height: 100%;
            background: linear-gradient(135deg, #00b09b 0%, #96c93d 100%);
            width: 0%;
            transition: width 0.3s ease;
        }
        
        .status-text {
            text-align: center;
            margin-top: 10px;
            color: #666;
            font-size: 0.9em;
        }
        
        .features {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin-top: 30px;
        }
        
        .feature {
            text-align: center;
            padding: 20px;
            background: #f8f9fa;
            border-radius: 8px;
        }
        
        .feature-icon {
            font-size: 2em;
            color: #4facfe;
            margin-bottom: 10px;
        }
        
        .footer {
            text-align: center;
            padding: 20px;
            color: #666;
            font-size: 0.9em;
            border-top: 1px solid #eee;
        }
        
        @media (max-width: 600px) {
            .container {
                margin: 10px;
            }
            
            .content {
                padding: 20px;
            }
            
            .upload-area {
                padding: 20px;
            }
            
            .header h1 {
                font-size: 1.8em;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎵 KGM → FLAC 转换器</h1>
            <p>安全、快速地将加密音频转换为标准FLAC格式</p>
        </div>
        
        <div class="content">
            <form id="uploadForm" action="/api/convert" method="post" enctype="multipart/form-data">
                <div class="upload-area" id="dropZone">
                    <div class="upload-icon">📁</div>
                    <h3>选择或拖放文件到此区域</h3>
                    <p>支持 .kgm, .kgma, .vpr 格式文件</p>
                    <p>最多可上传 {{.MaxFiles}} 个文件，单个文件不超过 {{.MaxFileSizeGB}} GB</p>
                    
                    <input type="file" id="fileInput" name="files" multiple 
                           accept=".kgm,.kgma,.vpr" class="file-input" />
                    <button type="button" class="browse-btn" onclick="document.getElementById('fileInput').click()">
                        选择文件
                    </button>
                </div>
                
                <div id="fileList" class="file-list"></div>
                
                <div class="progress-container" style="display: none;" id="progressContainer">
                    <div class="progress-bar">
                        <div class="progress-fill" id="progressFill"></div>
                    </div>
                    <div class="status-text" id="statusText">准备上传...</div>
                </div>
                
                <button type="submit" class="submit-btn" id="submitBtn" disabled>
                    开始转换
                </button>
            </form>
            
            <div class="features">
                <div class="feature">
                    <div class="feature-icon">🔒</div>
                    <h4>安全解密</h4>
                    <p>纯Go实现，无数据泄露风险</p>
                </div>
                <div class="feature">
                    <div class="feature-icon">⚡</div>
                    <h4>快速转换</h4>
                    <p>支持批量处理，高效转换</p>
                </div>
                <div class="feature">
                    <div class="feature-icon">🎧</div>
                    <h4>高质量输出</h4>
                    <p>转换为标准FLAC格式</p>
                </div>
            </div>
        </div>
        
        <div class="footer">
            <p>© 2025 KGM to FLAC Converter | 后端版本 {{.Version}}</p>
        </div>
    </div>

    <script>
        const maxFiles = {{.MaxFiles}};
        const maxFileSize = {{.MaxFileSize}};
        const maxFileSizeMB = {{.MaxFileSizeMB}};
        
        const dropZone = document.getElementById('dropZone');
        const fileInput = document.getElementById('fileInput');
        const fileList = document.getElementById('fileList');
        const submitBtn = document.getElementById('submitBtn');
        const progressContainer = document.getElementById('progressContainer');
        const progressFill = document.getElementById('progressFill');
        const statusText = document.getElementById('statusText');
        
        let selectedFiles = [];
        
        // 拖放功能
        ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
            dropZone.addEventListener(eventName, preventDefaults, false);
        });
        
        function preventDefaults(e) {
            e.preventDefault();
            e.stopPropagation();
        }
        
        ['dragenter', 'dragover'].forEach(eventName => {
            dropZone.addEventListener(eventName, highlight, false);
        });
        
        ['dragleave', 'drop'].forEach(eventName => {
            dropZone.addEventListener(eventName, unhighlight, false);
        });
        
        function highlight() {
            dropZone.classList.add('dragover');
        }
        
        function unhighlight() {
            dropZone.classList.remove('dragover');
        }
        
        dropZone.addEventListener('drop', handleDrop, false);
        
        function handleDrop(e) {
            const dt = e.dataTransfer;
            const files = dt.files;
            handleFiles(files);
        }
        
        // 文件选择处理
        fileInput.addEventListener('change', function(e) {
            handleFiles(e.target.files);
        });
        
        function handleFiles(files) {
            const newFiles = Array.from(files).filter(file => {
                const ext = file.name.toLowerCase().split('.').pop();
                return ['kgm', 'kgma', 'vpr'].includes(ext);
            });
            
            if (selectedFiles.length + newFiles.length > maxFiles) {
                alert('最多只能选择 ' + maxFiles + ' 个文件');
                return;
            }
            
            newFiles.forEach(file => {
                if (file.size > maxFileSize) {
                    alert('文件 ' + file.name + ' 超过 ' + maxFileSizeMB + ' MB 限制');
                    return;
                }
                
                selectedFiles.push(file);
                addFileToList(file);
            });
            
            updateSubmitButton();
        }
        
        function addFileToList(file) {
            const fileItem = document.createElement('div');
            fileItem.className = 'file-item';
            fileItem.innerHTML = '<div class="file-name">' + file.name + '</div>' +
                                '<div class="file-size">' + formatFileSize(file.size) + '</div>' +
                                '<button type="button" class="remove-btn" onclick="removeFile(\'' + file.name + '\')">移除</button>';
            fileList.appendChild(fileItem);
        }
        
        function removeFile(fileName) {
            selectedFiles = selectedFiles.filter(file => file.name !== fileName);
            renderFileList();
            updateSubmitButton();
        }
        
        function renderFileList() {
            fileList.innerHTML = '';
            selectedFiles.forEach(file => addFileToList(file));
        }
        
        function updateSubmitButton() {
            submitBtn.disabled = selectedFiles.length === 0;
        }
        
        function formatFileSize(bytes) {
            if (bytes === 0) return '0 Bytes';
            const k = 1024;
            const sizes = ['Bytes', 'KB', 'MB', 'GB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
        }
        
        // 表单提交处理
        document.getElementById('uploadForm').addEventListener('submit', function(e) {
            e.preventDefault();
            
            if (selectedFiles.length === 0) {
                alert('请至少选择一个文件');
                return;
            }
            
            const formData = new FormData();
            selectedFiles.forEach(file => {
                formData.append('files', file);
            });
            
            // 显示进度条
            progressContainer.style.display = 'block';
            submitBtn.disabled = true;
            statusText.textContent = '上传中...';
            
            fetch('/api/convert', {
                method: 'POST',
                body: formData
            })
            .then(response => {
                if (!response.ok) {
                    throw new Error('上传失败');
                }
                return response.blob();
            })
            .then(blob => {
                // 创建下载链接
                const url = window.URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.style.display = 'none';
                a.href = url;
                
                if (selectedFiles.length === 1) {
                    a.download = selectedFiles[0].name.replace(/\.[^/.]+$/, "") + '.flac';
                } else {
                    a.download = 'kgm2flac_result.zip';
                }
                
                document.body.appendChild(a);
                a.click();
                window.URL.revokeObjectURL(url);
                
                statusText.textContent = '转换完成！';
                progressFill.style.width = '100%';
                
                // 重置表单
                setTimeout(() => {
                    selectedFiles = [];
                    fileList.innerHTML = '';
                    updateSubmitButton();
                    progressContainer.style.display = 'none';
                    progressFill.style.width = '0%';
                }, 2000);
            })
            .catch(error => {
                console.error('Error:', error);
                statusText.textContent = '上传失败: ' + error.message;
                progressFill.style.width = '0%';
                submitBtn.disabled = false;
            });
        });
    </script>
</body>
</html>`

type ConvertHandler struct {
	cfg            *config.Config
	decryptService *service.DecryptService
}

func NewConvertHandler(cfg *config.Config) *ConvertHandler {
	return &ConvertHandler{
		cfg:            cfg,
		decryptService: service.NewDecryptService(),
	}
}

func (h *ConvertHandler) HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 准备模板数据
	templateData := map[string]interface{}{
		"MaxFiles":      h.cfg.MaxFiles,
		"MaxFileSize":   h.cfg.MaxFileSize,
		"MaxFileSizeGB": h.cfg.MaxFileSize >> 30,
		"MaxFileSizeMB": h.cfg.MaxFileSize >> 20,
		"Version":       "1.0.0",
	}

	t := template.Must(template.New("index").Parse(page))
	if err := t.Execute(w, templateData); err != nil {
		http.Error(w, "模板渲染失败", http.StatusInternalServerError)
		log.Printf("[ERR] template execute failed: %v", err)
	}
}

func (h *ConvertHandler) HandleConvert(w http.ResponseWriter, r *http.Request) {
	startReq := time.Now()
	clientIP := getClientIP(r)

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 限制整个请求体最大值
	limit := int64(h.cfg.MaxFiles)*h.cfg.MaxFileSize + (10 << 20) // +10MiB
	r.Body = http.MaxBytesReader(w, r.Body, limit)

	// ParseMultipartForm
	if err := r.ParseMultipartForm(h.cfg.ParseFormMemory); err != nil {
		http.Error(w, "表单解析失败: "+err.Error(), http.StatusBadRequest)
		log.Printf("[ERR] parse multipart form failed ip=%s err=%v", clientIP, err)
		return
	}
	defer func() {
		// 清理 ParseMultipartForm 创建的临时文件
		if r.MultipartForm != nil {
			_ = r.MultipartForm.RemoveAll()
		}
	}()

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "未选择文件（字段名为 files）", http.StatusBadRequest)
		return
	}
	if len(files) > h.cfg.MaxFiles {
		http.Error(w, fmt.Sprintf("最多上传 %d 个文件", h.cfg.MaxFiles), http.StatusBadRequest)
		return
	}

	log.Printf("[UPLOAD START] ip=%s files=%d", clientIP, len(files))

	// 创建临时工作目录
	workDir, err := os.MkdirTemp("", "kgm2flac_*")
	if err != nil {
		http.Error(w, "无法创建临时工作目录: "+err.Error(), http.StatusInternalServerError)
		log.Printf("[ERR] mkdir temp failed ip=%s err=%v", clientIP, err)
		return
	}
	defer func() {
		_ = os.RemoveAll(workDir)
	}()

	// 处理每个文件
	results := make([]types.ConvertResult, 0, len(files))
	for _, fh := range files {
		result := h.processSingleFile(fh, workDir, clientIP, r.Context())
		results = append(results, result)
	}

	// 统计成功数量
	successCount := 0
	for _, r := range results {
		if r.Err == nil {
			successCount++
		}
	}

	// 处理响应
	if successCount == 0 {
		http.Error(w, "所有文件处理失败", http.StatusBadRequest)
		return
	}

	if successCount == 1 {
		h.serveSingleFile(w, r, results, clientIP)
	} else {
		h.serveZipFile(w, r, results, workDir, clientIP)
	}

	totalDur := time.Since(startReq)
	log.Printf("[UPLOAD END] ip=%s total_files=%d success=%d took=%s", clientIP, len(files), successCount, totalDur)

	// 记录每个文件的详情
	for _, rr := range results {
		if rr.Err != nil {
			log.Printf("[FILE RESULT] ip=%s name=%s size=%d err=%v", clientIP, rr.OrigName, rr.Size, rr.Err)
		} else {
			log.Printf("[FILE RESULT] ip=%s name=%s size=%d out=%s dur=%s", clientIP, rr.OrigName, rr.Size, rr.OutPath, rr.Duration)
		}
	}
}

func (h *ConvertHandler) processSingleFile(fh *multipart.FileHeader, workDir, clientIP string, ctx interface{}) types.ConvertResult {
	start := time.Now()
	result := types.ConvertResult{
		OrigName: fh.Filename,
		Size:     fh.Size,
	}

	log.Printf("[FILE] ip=%s filename=%s size=%d", clientIP, fh.Filename, fh.Size)

	// 检查文件大小
	if fh.Size > h.cfg.MaxFileSize {
		result.Err = fmt.Errorf("文件 %s 超过单文件限制 (%d bytes)", fh.Filename, h.cfg.MaxFileSize)
		log.Printf("[ERR] %v", result.Err)
		return result
	}

	// 打开上传的文件
	f, err := fh.Open()
	if err != nil {
		result.Err = fmt.Errorf("打开上传文件失败: %w", err)
		log.Printf("[ERR] open uploaded file failed ip=%s name=%s err=%v", clientIP, fh.Filename, err)
		return result
	}
	defer f.Close()

	// 保存上传文件到临时位置
	inPath, cleanupIn, err := h.persistUpload(f, fh)
	if err != nil {
		result.Err = fmt.Errorf("保存上传文件失败: %w", err)
		log.Printf("[ERR] persist upload failed ip=%s name=%s err=%v", clientIP, fh.Filename, err)
		return result
	}
	defer cleanupIn()

	// 解密文件
	outRaw, cleanupRaw, err := h.decryptService.DecryptKgmFile(inPath)
	if err != nil {
		result.Err = fmt.Errorf("解密失败: %w", err)
		log.Printf("[ERR] decrypt failed ip=%s name=%s err=%v", clientIP, fh.Filename, err)
		return result
	}
	defer cleanupRaw()

	// 嗅探音频格式
	rawExt, err := h.sniffAudioExt(outRaw)
	if err != nil {
		result.Err = fmt.Errorf("识别音频格式失败: %w", err)
		log.Printf("[ERR] sniff audio ext failed ip=%s name=%s err=%v", clientIP, fh.Filename, err)
		return result
	}

	// 处理输出文件
	finalPath := filepath.Join(workDir, utils.ReplaceExt(fh.Filename, ".flac"))
	if rawExt == ".flac" {
		// 如果已经是flac，直接重命名
		if err := os.Rename(outRaw, finalPath); err != nil {
			if err := h.copyFile(outRaw, finalPath); err != nil {
				result.Err = fmt.Errorf("移动FLAC文件失败: %w", err)
				log.Printf("[ERR] move/copy flac failed ip=%s name=%s err=%v", clientIP, fh.Filename, err)
				return result
			}
			_ = os.Remove(outRaw)
		}
	} else {
		// 需要转码为FLAC
		if err := h.convertToFlac(outRaw, finalPath); err != nil {
			result.Err = fmt.Errorf("转码为FLAC失败: %w", err)
			log.Printf("[ERR] ffmpeg convert failed ip=%s name=%s err=%v", clientIP, fh.Filename, err)
			return result
		}
		_ = os.Remove(outRaw)
	}

	result.OutPath = finalPath
	result.Duration = time.Since(start)
	log.Printf("[FILE DONE] ip=%s name=%s out=%s dur=%s", clientIP, fh.Filename, finalPath, result.Duration)

	return result
}

func (h *ConvertHandler) serveSingleFile(w http.ResponseWriter, r *http.Request, results []types.ConvertResult, clientIP string) {
	var fileToServe string
	var origName string

	for _, rr := range results {
		if rr.Err == nil {
			fileToServe = rr.OutPath
			origName = rr.OrigName
			break
		}
	}

	if fileToServe == "" {
		http.Error(w, "内部错误：没有可下载的文件", http.StatusInternalServerError)
		return
	}

	log.Printf("[RESP] ip=%s serve single file=%s size=%d", clientIP, fileToServe, utils.FileSizeSafe(fileToServe))
	w.Header().Set("Content-Type", "audio/flac")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", utils.ReplaceExt(origName, ".flac")))
	http.ServeFile(w, r, fileToServe)
}

func (h *ConvertHandler) serveZipFile(w http.ResponseWriter, r *http.Request, results []types.ConvertResult, workDir, clientIP string) {
	zipPath := filepath.Join(workDir, "kgm2flac_result_"+utils.RandHex(8)+".zip")
	zipFile, err := os.Create(zipPath)
	if err != nil {
		http.Error(w, "无法创建zip文件: "+err.Error(), http.StatusInternalServerError)
		log.Printf("[ERR] create zip failed ip=%s err=%v", clientIP, err)
		return
	}
	defer zipFile.Close()

	zw := zip.NewWriter(zipFile)
	defer zw.Close()

	// 添加成功文件到zip
	successCount := 0
	for _, rr := range results {
		if rr.Err != nil || rr.OutPath == "" {
			continue
		}
		if err := h.addFileToZip(zw, rr.OutPath, filepath.Base(rr.OutPath)); err != nil {
			log.Printf("[ERR] add to zip failed ip=%s file=%s err=%v", clientIP, rr.OutPath, err)
			continue
		}
		successCount++
	}

	if err := zw.Close(); err != nil {
		http.Error(w, "无法生成zip: "+err.Error(), http.StatusInternalServerError)
		log.Printf("[ERR] close zip failed ip=%s err=%v", clientIP, err)
		return
	}

	if successCount == 0 {
		http.Error(w, "没有成功的文件可打包", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="kgm2flac_result.zip"`)
	http.ServeFile(w, r, zipPath)
}

func (h *ConvertHandler) persistUpload(src multipart.File, hdr *multipart.FileHeader) (path string, cleanup func(), err error) {
	// 读取开头4字节以便后续检查
	b := make([]byte, 4)
	if _, err := io.ReadFull(src, b); err != nil && err != io.EOF {
		return "", func() {}, err
	}
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		// 忽略Seek错误，继续处理
	}

	name := fmt.Sprintf("kgm_%s%s", utils.RandHex(8), filepath.Ext(hdr.Filename))
	path = filepath.Join(os.TempDir(), name)
	f, err := os.Create(path)
	if err != nil {
		return "", func() {}, err
	}
	defer f.Close()

	if _, err = io.Copy(f, src); err != nil {
		return "", func() {}, err
	}

	return path, func() { _ = os.Remove(path) }, nil
}

func (h *ConvertHandler) sniffAudioExt(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	head := make([]byte, 12)
	if _, err := io.ReadFull(f, head); err != nil {
		return "", err
	}

	switch {
	case bytes.HasPrefix(head, []byte("fLaC")):
		return ".flac", nil
	case bytes.HasPrefix(head, []byte("ID3")):
		return ".mp3", nil
	case head[0] == 0xFF && (head[1]&0xE0) == 0xE0:
		return ".mp3", nil
	case bytes.HasPrefix(head, []byte("OggS")):
		return ".ogg", nil
	default:
		return "", fmt.Errorf("未知音频头: %x", head)
	}
}

func (h *ConvertHandler) convertToFlac(inputPath, outputPath string) error {
	cmd := exec.Command(h.cfg.FFmpegBin,
		"-y",
		"-hide_banner",
		"-loglevel", "error",
		"-i", inputPath,
		"-map_metadata", "-1",
		outputPath,
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg执行失败: %w", err)
	}
	return nil
}

func (h *ConvertHandler) copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

func (h *ConvertHandler) addFileToZip(zw *zip.Writer, path string, nameInZip string) error {
	finfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	fh, err := zip.FileInfoHeader(finfo)
	if err != nil {
		return err
	}

	fh.Name = nameInZip
	fh.Method = zip.Deflate

	w, err := zw.CreateHeader(fh)
	if err != nil {
		return err
	}

	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()

	_, err = io.Copy(w, in)
	return err
}

// StartServer 启动HTTP服务器
func StartServer(cfg *config.Config) error {
	handler := NewConvertHandler(cfg)
	mux := http.NewServeMux()

	mux.HandleFunc("/", handler.HandleRoot)
	mux.HandleFunc("/api/convert", handler.HandleConvert)

	log.Printf("启动服务器，监听地址: %s", cfg.Addr)
	log.Printf("FFmpeg路径: %s", cfg.FFmpegBin)
	log.Printf("单文件最大大小: %d bytes", cfg.MaxFileSize)
	log.Printf("最大文件数: %d", cfg.MaxFiles)

	return http.ListenAndServe(cfg.Addr, logRequest(mux))
}
