package global

import (
	"errors"
)

var (
	ERR_ROOT_NODE_NOT_EXIST = errors.New("Root node is not exist")
	ERR_NODE_NOT_EXIST      = errors.New("The node is not exist")
	ERR_CACHE_KEY_NOT_EXIST = errors.New("The key is not exist")
	ERR_RE_INDEX            = errors.New("The data is running reindex.")
	ERR_NODE_LEVEL_IS_ZERO  = errors.New("Node 's level must be not  zero.")
	ERR_PERMISSION_DENY     = errors.New("Permission deny.")
)
