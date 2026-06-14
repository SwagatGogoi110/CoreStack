package model

import (
	"time"
)

type User struct {
	UserID                 string            `json:"userId"`
	UserName               string            `json:"userName"`
	Path                   string            `json:"path"`
	Arn                    string            `json:"arn"`
	CreateDate             time.Time         `json:"createDate"`
	PasswordLastUsed       *time.Time        `json:"passwordLastUsed,omitempty"`
	Tags                   map[string]string `json:"tags"`
	GroupNames             []string          `json:"groupNames"`
	AttachedPolicyArns     []string          `json:"attachedPolicyArns"`
	InlinePolicies         map[string]string `json:"inlinePolicies"`
	PermissionsBoundaryArn string            `json:"permissionsBoundaryArn,omitempty"`
}

type AccessKey struct {
	AccessKeyID     string    `json:"accessKeyId"`
	SecretAccessKey string    `json:"secretAccessKey"`
	UserName        string    `json:"userName"`
	Status          string    `json:"status"` // Active | Inactive
	CreateDate      time.Time `json:"createDate"`
}

type Role struct {
	RoleID                   string            `json:"roleId"`
	RoleName                 string            `json:"roleName"`
	Path                     string            `json:"path"`
	Arn                      string            `json:"arn"`
	AssumeRolePolicyDocument string            `json:"assumeRolePolicyDocument"`
	Description              string            `json:"description"`
	MaxSessionDuration       int               `json:"maxSessionDuration"`
	CreateDate               time.Time         `json:"createDate"`
	Tags                     map[string]string `json:"tags"`
	AttachedPolicyArns       []string          `json:"attachedPolicyArns"`
	InlinePolicies           map[string]string `json:"inlinePolicies"`
	PermissionsBoundaryArn   string            `json:"permissionsBoundaryArn,omitempty"`
}

type Policy struct {
	PolicyID         string                   `json:"policyId"`
	PolicyName       string                   `json:"policyName"`
	Path             string                   `json:"path"`
	Arn              string                   `json:"arn"`
	Description      string                   `json:"description"`
	DefaultVersionID string                   `json:"defaultVersionId"`
	AttachmentCount  int                      `json:"attachmentCount"`
	CreateDate       time.Time                `json:"createDate"`
	UpdateDate       time.Time                `json:"updateDate"`
	Tags             map[string]string        `json:"tags"`
	Versions         map[string]PolicyVersion `json:"versions"`
}

type PolicyVersion struct {
	VersionID        string    `json:"versionId"`
	Document         string    `json:"document"`
	IsDefaultVersion bool      `json:"isDefaultVersion"`
	CreateDate       time.Time `json:"createDate"`
}

type Group struct {
	GroupID            string            `json:"groupId"`
	GroupName          string            `json:"groupName"`
	Path               string            `json:"path"`
	Arn                string            `json:"arn"`
	CreateDate         time.Time         `json:"createDate"`
	AttachedPolicyArns []string          `json:"attachedPolicyArns"`
	InlinePolicies     map[string]string `json:"inlinePolicies"`
}
