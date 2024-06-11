package fakes3access

type HashMethod string

const (
	Blake2b HashMethod = "blake2b"
	MD5     HashMethod = "md5"
	SHA1    HashMethod = "sha1"
	SHA256  HashMethod = "sha256"
	SHA512  HashMethod = "sha512"
)

// 一个服务器的任何配置文件信息, 我会尽可能存储到配置文件而非数据库中, 考虑到数据库并不太适合存储配置信息.

//
//
// bucket.json也同时定义了Hash方法和位数。也方便重构数据库时采用的策略。包含以下属性：
//
