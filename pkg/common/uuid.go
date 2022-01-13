package common

import "github.com/google/uuid"

/*
 v1,v4都是每次生成一个唯一的ID，而v1同一时刻的输出非常相似，v1末尾nodeID部分用的都是mac地址，前面time的mid,high以及clock序列都是一样的，
 只有time-low部分不同。v4加入了随机数，对各个部分都进行了随机处理，同一时刻的输出差别很大。

 v2 NewDCEGroup()根据os.Getgid取到的用户组ID来生成uuid,同一时刻的输出是相同的。

 v2 NewDCEPerson()根据os.Getuid取到的用户ID来生成uuid,同一时刻的输出也是相同的。

 v3 NewMD5(space UUID, data []byte)是根据参数传入的UUID结构体和[]byte再重新转换一次。只要传入参数相同则任意时刻的输出也相同。

 v5NewSHA1(space UUID, data []byte)是根据参数传入的UUID结构体和[]byte再重新转换一次。只要传入参数相同则任意时刻的输出也相同。
*/

func UuidV4() string {
	id := uuid.New()
	return id.String()
}

func UuidV1() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func UuidV2G() (string, error) {
	id, err := uuid.NewDCEGroup()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func UuidV2P() (string, error) {
	id, err := uuid.NewDCEPerson()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func UuidV3() (string, error) {
	tmp, err := uuid.NewDCEPerson()
	if err != nil {
		return "", err
	}
	id := uuid.NewMD5(tmp, []byte("fssds32"))
	return id.String(), nil
}

func UuidV5() (string, error) {
	tmp, err := uuid.NewDCEPerson()
	if err != nil {
		return "", err
	}
	id := uuid.NewSHA1(tmp, []byte("fssds32"))
	return id.String(), nil
}
