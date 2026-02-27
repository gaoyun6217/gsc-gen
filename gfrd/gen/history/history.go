package history

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// GenerationRecord 生成记录
type GenerationRecord struct {
	ID           string            `json:"id"`           // 记录 ID
	Table        string            `json:"table"`        // 表名
	Module       string            `json:"module"`       // 模块名
	GeneratedAt  time.Time         `json:"generated_at"` // 生成时间
	Files        []GeneratedFile   `json:"files"`        // 生成的文件列表
	TableComment string            `json:"table_comment"`// 表注释
	FieldCount   int               `json:"field_count"`  // 字段数量
	Config       GeneratorConfig   `json:"config"`       // 生成配置
	Checksum     string            `json:"checksum"`     // 文件校验和
}

// GeneratedFile 生成的文件
type GeneratedFile struct {
	Path      string    `json:"path"`       // 文件路径
	Type      string    `json:"type"`       // 文件类型：backend/frontend/sql
	Content   string    `json:"content"`    // 文件内容（用于回滚）
	Checksum  string    `json:"checksum"`   // 文件校验和
	CreatedAt time.Time `json:"created_at"` // 创建时间
}

// GeneratorConfig 生成器配置快照
type GeneratorConfig struct {
	Output    string `json:"output"`    // 后端输出目录
	WebOutput string `json:"web_output"`// 前端输出目录
	Package   string `json:"package"`   // Go 包名
	Features  string `json:"features"`  // 功能列表
}

// HistoryManager 历史记录管理器
type HistoryManager struct {
	historyDir string
	records    []*GenerationRecord
}

// NewHistoryManager 创建历史记录管理器
func NewHistoryManager(historyDir string) (*HistoryManager, error) {
	hm := &HistoryManager{
		historyDir: historyDir,
		records:    make([]*GenerationRecord, 0),
	}

	// 创建历史目录
	if err := os.MkdirAll(historyDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create history directory: %w", err)
	}

	// 加载历史记录
	if err := hm.LoadRecords(); err != nil {
		return nil, fmt.Errorf("failed to load records: %w", err)
	}

	return hm, nil
}

// LoadRecords 加载历史记录
func (hm *HistoryManager) LoadRecords() error {
	historyFile := filepath.Join(hm.historyDir, "history.json")

	data, err := os.ReadFile(historyFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // 文件不存在，返回空记录
		}
		return err
	}

	return json.Unmarshal(data, &hm.records)
}

// SaveRecords 保存历史记录
func (hm *HistoryManager) SaveRecords() error {
	historyFile := filepath.Join(hm.historyDir, "history.json")

	data, err := json.MarshalIndent(hm.records, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(historyFile, data, 0644)
}

// AddRecord 添加生成记录
func (hm *HistoryManager) AddRecord(record *GenerationRecord) error {
	hm.records = append(hm.records, record)

	// 按时间排序
	sort.Slice(hm.records, func(i, j int) bool {
		return hm.records[i].GeneratedAt.Before(hm.records[j].GeneratedAt)
	})

	return hm.SaveRecords()
}

// GetRecords 获取所有记录
func (hm *HistoryManager) GetRecords() []*GenerationRecord {
	return hm.records
}

// GetRecordsByTable 根据表名获取记录
func (hm *HistoryManager) GetRecordsByTable(tableName string) []*GenerationRecord {
	var result []*GenerationRecord
	for _, r := range hm.records {
		if r.Table == tableName {
			result = append(result, r)
		}
	}
	return result
}

// GetLatestRecord 获取最新记录
func (hm *HistoryManager) GetLatestRecord() *GenerationRecord {
	if len(hm.records) == 0 {
		return nil
	}
	return hm.records[len(hm.records)-1]
}

// Rollback 回滚到指定记录
func (hm *HistoryManager) Rollback(recordID string) error {
	record := hm.GetRecordByID(recordID)
	if record == nil {
		return fmt.Errorf("record not found: %s", recordID)
	}

	// 恢复每个文件
	for _, file := range record.Files {
		if err := os.WriteFile(file.Path, []byte(file.Content), 0644); err != nil {
			return fmt.Errorf("failed to restore file %s: %w", file.Path, err)
		}
	}

	return nil
}

// GetRecordByID 根据 ID 获取记录
func (hm *HistoryManager) GetRecordByID(id string) *GenerationRecord {
	for _, r := range hm.records {
		if r.ID == id {
			return r
		}
	}
	return nil
}

// DeleteRecord 删除记录
func (hm *HistoryManager) DeleteRecord(recordID string) error {
	for i, r := range hm.records {
		if r.ID == recordID {
			hm.records = append(hm.records[:i], hm.records[i+1:]...)
			return hm.SaveRecords()
		}
	}
	return fmt.Errorf("record not found: %s", recordID)
}

// ClearHistory 清空历史记录
func (hm *HistoryManager) ClearHistory() error {
	hm.records = make([]*GenerationRecord, 0)
	return hm.SaveRecords()
}

// GenerateRecordID 生成记录 ID
func GenerateRecordID() string {
	return fmt.Sprintf("gen_%d", time.Now().UnixNano())
}

// CalculateChecksum 计算文件校验和
func CalculateChecksum(content string) string {
	hash := md5.Sum([]byte(content))
	return hex.EncodeToString(hash[:])
}
