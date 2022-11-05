import http from 'k6/http';

export default function () {
  const url = 'http://localhost:8080/transactions';
  const payload = JSON.stringify({
	"items": [
		{
			"quantity": 2,
			"item_id": "9fcb59b5-9b2a-4759-abb7-6b2809cb47e8"
		},
		{
			"quantity": -2,
			"item_id": "9fcb59b5-9b2a-4759-abb7-6b2809cb47e8"
		},
		{
			"quantity": 1,
			"item_id": "557adb58-8b4d-4d06-85ed-ced7b5ff22fc"
		}		
	],
	"location_id": "5996ff4d-725d-42c8-953e-99b47412db14"
  });

  const params = {
    headers: {
      'Authorization': 'Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IkhqNFdyLXpJa05XTFdTUnVQZ2ZlZCJ9.eyJpc3MiOiJodHRwczovL3dhcmVob3VzZS5hdS5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NjM1MzQ4YWI5ZmRmM2IxNWI1ZGY1NDE3IiwiYXVkIjpbImh0dHA6Ly9sb2NhbGhvc3Q6ODA4MC8iLCJodHRwczovL3dhcmVob3VzZS5hdS5hdXRoMC5jb20vdXNlcmluZm8iXSwiaWF0IjoxNjY3NjA3MzQzLCJleHAiOjE2Njc2OTM3NDMsImF6cCI6ImlyRUkxZFZKT2J2aUo2N21ac1lGRWpQOTNFSWx0aHlzIiwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCB3cml0ZTpjYXRlZ29yeSByZWFkOmNhdGVnb3J5IHdyaXRlOml0ZW0gcmVhZDppdGVtIG9mZmxpbmVfYWNjZXNzIn0.lPcszwxZsUn0NGkRilzCjMfTcmXFqrHSCM_TrFX_22_2_sakmnxQAtz4ijpo-Xlor45UzxDaHOrVZnZHPitxCdMDJxN3_NysMlUcF35O1gZkhNcU8Kh4uRsVVsok36DcIvvYHUwO-N2zhgYxN1Ulku5jzwsUlcZi4AbEPwtzrn55OLyZdBej9t0BC91N4gnTc5v7I0S98dFuhH0YkY-qkXE03KSWOgZn5MX8ywsd4B6xaHpqvXr27kmldg-qdW8T7wQGRGhtouCV763PvrjIbe6yL4ZkDfiSKF6gB6iRbrGe9DGDEv7FleTlpOIxlCaD1O_CyX9ccjMiiaYQq5Mz3Q',
      'Content-Type': 'application/json',
    },
  };

  http.post(url, payload, params);
}
