package email

import (
	"e-ticketing-gin/configs"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"math/rand"
)

type EmailInterface interface {
	SendEmail(to, subject, body string) error
	HTMLBodyReset(username string) (string, string, string)
	HTMLBodyVerification(username string) (string, string, string)
}

type Email struct {
	c *configs.ProgramConfig
}

func NewEmail(c *configs.ProgramConfig) EmailInterface {
	return &Email{
		c: c,
	}
}

func (e *Email) SendEmail(to, subject, body string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", e.c.Email)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, e.c.Email, e.c.Password)

	err := dialer.DialAndSend(message)
	if err != nil {
		logrus.Error("ERROR : Dialer Erorr : ", err.Error())
		return err
	}

	return nil
}

func (e *Email) generateRandomCode(length int) string {
	const charset = "0123456789"
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}

	return string(code)
}

func (e *Email) HTMLBodyReset(username string) (string, string, string) {
	code := e.generateRandomCode(4)
	header, htmlBody := e.htmlBodyEmailReset(username, code)

	return header, htmlBody, code
}

func (e *Email) HTMLBodyVerification(username string) (string, string, string) {
	code := e.generateRandomCode(4)
	header, htmlBody := e.htmlBodyEmailVerification(username, code)

	return header, htmlBody, code
}

func (e *Email) htmlBodyEmailReset(username, code string) (string, string) {
	header := "Pemulihan Kata Sandi - Kode OTP Dikirimkan untuk Anda"
	htmlBody := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta http-equiv="X-UA-Compatible" content="IE=edge">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Kode Verifikasi</title>
		</head>
		<body style="margin: 0; padding: 0; box-sizing: border-box;">
			<table align="center" cellpadding="0" cellspacing="0" width="95%">
			<tr>
				<td align="center">
				<table align="center" cellpadding="0" cellspacing="0" width="600" style="border-spacing: 2px 5px;" bgcolor="#fff">
				<tr>
                        <td style="background-color: #fff; text-align: center; padding: 20px;">
                            <img src="https://i.ibb.co.com/3RZSKjL/Golang-Email-Header.png" alt="Logo" style="width: 700px; height: auto;">
                        </td>
                    </tr>	
					<tr>
					<td bgcolor="#fff">
						<table cellpadding="0" cellspacing="0" width="100%%">
						<tr>
							<td style="padding: 10px 0 10px 0; font-family: Nunito, sans-serif; font-size: 20px; font-weight: 900">
							Halo, ` + username + `
							</td>
						</tr>
						</table>
					</td>
					</tr>
					<tr>
					<td bgcolor="#fff">
						<table cellpadding="0" cellspacing="0" width="100%%">
						<tr>
							<td style="padding: 0; font-family: Nunito, sans-serif; font-size: 16px;">
							Kami melihat bahwa Anda mengalami kesulitan untuk mengakses akun Anda. Jangan khawatir, kami di sini untuk membantu Anda! Kami telah mengirimkan kode OTP ke alamat email terkait dengan akun Anda.
							</td>
						</tr>
						<tr>
							<td style="padding: 20px 0 20px 0; font-family: Nunito, sans-serif; font-size: 16px; text-align: center;">
							<button style="background-color: #0085FF; border: none; color: white; padding: 15px 30px; text-align: center; display: inline-block; font-family: Nunito, sans-serif; font-size: 20px; font-weight: bold; cursor: pointer; border-radius:8px; letter-spacing: 10px;">
								` + code + `
							</button>
							</td>
						</tr>
						<tr>
							<td style="padding: 0; font-family: Nunito, sans-serif; font-size: 16px;">
							Silakan gunakan kode ini untuk mengatur ulang kata sandi Anda dengan mudah. Pastikan untuk segera mengganti kata sandi setelah berhasil masuk kembali ke akun Anda.
							<p></p>
							</td>
						</tr>
						<tr>
							<td style="padding: 0; font-family: Nunito, sans-serif; font-size: 16px;">
							Jika Anda tidak meminta pemulihan kata sandi ini, mohon abaikan pesan ini untuk menjaga keamanan akun Anda. 
							</td>
						</tr>
						</table>
					</td>
					</tr>
				</table>
				</td>
			</tr>
			</table>
		</body>
		</html>
		`

	return header, htmlBody
}

func (e *Email) htmlBodyEmailVerification(username, code string) (string, string) {
	header := "Verifikasi Akun anda - Kode OTP Dikirimkan untuk Anda"
	htmlBody := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta http-equiv="X-UA-Compatible" content="IE=edge">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Kode Verifikasi Akun</title>
		</head>
		<body style="margin: 0; padding: 0; box-sizing: border-box;">
			<table align="center" cellpadding="0" cellspacing="0" width="95%">
			<tr>
				<td align="center">
				<table align="center" cellpadding="0" cellspacing="0" width="600" style="border-spacing: 2px 5px;" bgcolor="#fff">
				<tr>
                        <td style="background-color: #fff; text-align: center; padding: 20px;">
                            <img src="https://i.ibb.co.com/3RZSKjL/Golang-Email-Header.png" alt="Logo" style="width: 700px; height: auto;">
                        </td>
                    </tr>	
					<tr>
					<td bgcolor="#fff">
						<table cellpadding="0" cellspacing="0" width="100%%">
						<tr>
							<td style="padding: 10px 0 10px 0; font-family: Nunito, sans-serif; font-size: 20px; font-weight: 900">
							Halo, ` + username + `
							</td>
						</tr>
						</table>
					</td>
					</tr>
					<tr>
					<td bgcolor="#fff">
						<table cellpadding="0" cellspacing="0" width="100%%">
						<tr>
							<td style="padding: 0; font-family: Nunito, sans-serif; font-size: 16px;">
							Terimakasih sudah mendaftar, berikut code verifikasi anda untuk mengaktifkan akun
							</td>
						</tr>
						<tr>
							<td style="padding: 20px 0 20px 0; font-family: Nunito, sans-serif; font-size: 16px; text-align: center;">
							<button style="background-color: #0085FF; border: none; color: white; padding: 15px 30px; text-align: center; display: inline-block; font-family: Nunito, sans-serif; font-size: 20px; font-weight: bold; cursor: pointer; border-radius:8px; letter-spacing: 10px;">
								` + code + `
							</button>
							</td>
						</tr>
						<tr>
							<td style="padding: 0; font-family: Nunito, sans-serif; font-size: 16px;">
							Silahkan gunakan kode ini untuk mengaktifkan akun ada, harap segera memasukan code sebelum waktu habis
							<p></p>
							</td>
						</tr>
						<tr>
							<td style="padding: 0; font-family: Nunito, sans-serif; font-size: 16px;">
							Jika Anda tidak meminta code verifikasi ini, mohon abaikan pesan ini untuk menjaga keamanan akun Anda. 
							</td>
						</tr>
						</table>
					</td>
					</tr>
				</table>
				</td>
			</tr>
			</table>
		</body>
		</html>
		`

	return header, htmlBody
}
