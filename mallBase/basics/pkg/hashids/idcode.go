package hashids

import (
	"strings"
)

// Hashids 封装hashids方法
type Hashids struct {
	HashID     *HashID
	HashIDData *HashIDData
	prefix     string
}

// NewHashids 创建Hashids对象
// salt可以使用用户创建记录时的用户唯一身份标识+当前时间戳的字符串作为值
// minLength指定转换后的最小长度,随着数字ID的增大长度可能会变长
func NewHashids(salt string, minLength int, prefix string) (*Hashids, error) {
	hd := NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	h, err := NewWithData(hd)
	if err != nil {
		return nil, err
	}
	return &Hashids{
		HashID:     h,
		HashIDData: hd,
		prefix:     prefix,
	}, nil
}

// Encode 将数字类型的ID转换为指定长度的随机字符串ID
// int64ID为需要转换的数字id，在没有自增主键ID时，可以采用当前用户已存在的记录数+1作为数字id，保证该数字在该用户下唯一
func (h *Hashids) Encode(int64ID int64) (string, error) {
	idstr, err := h.HashID.EncodeInt64([]int64{int64ID})
	if err != nil {
		return "", err
	}
	return h.prefix + idstr, nil
}

// Decode 将生成的随机字符串ID转为为原始的数字类型ID
func (h *Hashids) Decode(hashID string) (int64, error) {
	if h.prefix != "" {
		hashID = strings.TrimPrefix(hashID, h.prefix)
	}
	idSlice, err := h.HashID.DecodeInt64WithError(hashID)
	if err != nil {
		return 0, err
	}
	return idSlice[0], nil
}

// MustNewHashids 创建Hashids对象
// salt可以使用用户创建记录时的用户唯一身份标识+当前时间戳的字符串作为值
// minLength指定转换后的最小长度,随着数字ID的增大长度可能会变长
func MustNewHashids(salt string, minLength int, prefix string) *Hashids {
	hd := NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	h, err := NewWithData(hd)
	if err != nil {
		panic(err)
	}
	return &Hashids{
		HashID:     h,
		HashIDData: hd,
		prefix:     prefix,
	}
}

// IntHashEncode 返回int的hash值
func IntHashEncode(id int64, salt string, minLength int, prefix string) string {
	h, err := NewHashids(salt, minLength, prefix)
	if err != nil {
		return ""
	}
	s, err := h.Encode(id)
	if err != nil {
		return ""
	}
	return s
}

// IntHashDecode 返回hash值对应的int
func IntHashDecode(s string, salt string, minLength int, prefix string) int {
	h, err := NewHashids(salt, minLength, prefix)
	if err != nil {
		return 0
	}
	i, err := h.Decode(s)
	if err != nil {
		return 0
	}
	return int(i)
}
