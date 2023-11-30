package service

import (
	dto "GinProject/dto/user"
	"GinProject/utils"
)

func AddComment(comment *dto.CommentForm) bool {
	err := utils.Publish("amq.direct", "comment", comment)
	if err != nil {
		return false
	}
	return true
}
