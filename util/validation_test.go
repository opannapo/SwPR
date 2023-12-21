package util_test

import (
	"swpr/util"
	"testing"
)

func TestValidatePhoneForRegister(t *testing.T) {
	type args struct {
		phoneNumber string
	}
	tests := []struct {
		name       string
		args       args
		wantResult bool
	}{
		{
			name: "Success 10 number",
			args: args{
				phoneNumber: "+621234567890",
			},
			wantResult: true,
		},
		{
			name: "Success 12 number",
			args: args{
				phoneNumber: "+62123456789012",
			},
			wantResult: true,
		},
		{
			name: "Success 13 number",
			args: args{
				phoneNumber: "+621234567890123",
			},
			wantResult: true,
		},
		{
			name: "Error 9 number",
			args: args{
				phoneNumber: "+62123456789",
			},
			wantResult: false,
		},
		{
			name: "Error 14 number",
			args: args{
				phoneNumber: "+6212345678901234",
			},
			wantResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actualResult := util.ValidatePhoneForRegister(tt.args.phoneNumber); actualResult != tt.wantResult {
				t.Errorf("TestValidatePhoneForRegister() = %v, want %v", actualResult, tt.wantResult)
			}
		})
	}
}

func TestValidateNameForRegister(t *testing.T) {
	type args struct {
		phoneNumber string
	}
	tests := []struct {
		name       string
		args       args
		wantResult bool
	}{
		{
			name: "Success 3 character",
			args: args{
				phoneNumber: "opa",
			},
			wantResult: true,
		},
		{
			name: "Success 60 character",
			args: args{
				phoneNumber: "123456789012345678901234567890123456789012345678901234567890",
			},
			wantResult: true,
		},
		{
			name: "Error <3 character",
			args: args{
				phoneNumber: "ab",
			},
			wantResult: false,
		},
		{
			name: "Error >60 character",
			args: args{
				phoneNumber: "1234567890123456789012345678901234567890123456789012345678901",
			},
			wantResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actualResult := util.ValidateNameForRegister(tt.args.phoneNumber); actualResult != tt.wantResult {
				t.Errorf("ValidateNameForRegister() = %v, want %v", actualResult, tt.wantResult)
			}
		})
	}
}
