# bEDU (Bedah Edukasi)

Backend for Capstone Project SIB X Dicoding Cycle 4

## Daftar Isi

- [bEDU (Bedah Edukasi)](#bEDU)
  - [Daftar Isi](#daftar-isi)
  - [Daftar API](#daftar-api)
    - [Registration](#registration)
    - [Users](#users)

## Daftar API

Kumpulan API untuk bEDU (Bedah Edukasi) yang digunakan

### Registration

| Method | Developer | Endpoint  |                                        URL                                         |                                          Status                                          | Deskripsi           | Penggunaan `Authorization` |
| ------ | :-------: | :-------: | :--------------------------------------------------------------------------------: | :--------------------------------------------------------------------------------------: | :------------------ | :------------------------: |
| POST   |   Raha    |  /login   |  [Link](http://ec2-54-66-56-235.ap-southeast-2.compute.amazonaws.com:8080/login)   | `200` Ok `400` Bad Request `401` Unauthorized `404` Not found `500`Internal Server Error | Endpoint Login      |           Tidak            |
| POST   |   Raha    | /register | [Link](http://ec2-54-66-56-235.ap-southeast-2.compute.amazonaws.com:8080/register) | `201` Ok `400` Bad Request `401` Unauthorized `404` Not found `500`Internal Server Error | Endpoint Registrasi |           Tidak            |
