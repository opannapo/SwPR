package util_test

import (
	"swpr/util"
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	type args struct {
		password string
		hash     string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Success pwd 123456",
			args: args{
				password: "123456",
				hash:     "$2a$12$ya8fM2J7lPJrDkcFI.jIi.jpOtg3hBC7otnK7k4pt8svU2/RvGCam",
			},
			want: true,
		},
		{
			name: "Success pwd TestSawitPro!23$",
			args: args{
				password: "TestSawitPro!23$",
				hash:     "$2a$12$FjtOZ.O9SaaKs6XArCM40uA2Y4mVVXfwpZOE8SOyg/ssx51/YEv16",
			},
			want: true,
		},
		{
			name: "Error pwd TestSawitPro",
			args: args{
				password: "TestSawitPro",
				hash:     "$2a$12$FjtOZ.O9SaaKs6XArCM40uA2Y4mVVXfwpZOE8SOyg/ssx51/YEv16",
			},
			want: false,
		},
	}
	pwdUtil := util.NewPasswordUtil()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pwdUtil.CheckPasswordHash(tt.args.password, tt.args.hash); got != tt.want {
				t.Errorf("CheckPasswordHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name                 string
		args                 args
		isGeneratedHashValid bool
	}{
		{
			name: "Success",
			args: args{
				password: "123456",
			},
			isGeneratedHashValid: true,
		},
	}
	pwdUtil := util.NewPasswordUtil()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := pwdUtil.HashPassword(tt.args.password)
			actualMatch := pwdUtil.CheckPasswordHash(tt.args.password, got)
			if tt.isGeneratedHashValid != actualMatch {
				t.Errorf("HashPassword() got = %v, want isMatch %v, actualMatch %v", got, tt.isGeneratedHashValid, actualMatch)
			}
		})
	}
}
