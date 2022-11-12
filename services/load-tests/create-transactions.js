import http from 'k6/http';

export default function () {
  const url = 'http://localhost:8080/transactions';
  const payload = JSON.stringify({
	"items": [
		{
			"quantity": 2,
			"item_id": "abfdd2a2-cd55-4e51-b6eb-d3ae0dde8f8e",
		},
		{
			"quantity": 1,
			"item_id": "17b593bb-f0db-4707-861a-53a8a3f41e7e"
		}		
	],
	"location_id": "35f3c5d9-2950-4f58-b962-1557d374909e"
  });

  const params = {
    headers: {
      'Authorization': 'Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IkhqNFdyLXpJa05XTFdTUnVQZ2ZlZCJ9.eyJpc3MiOiJodHRwczovL3dhcmVob3VzZS5hdS5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NjM1MzQ4YWI5ZmRmM2IxNWI1ZGY1NDE3IiwiYXVkIjpbImh0dHA6Ly9sb2NhbGhvc3Q6ODA4MC8iLCJodHRwczovL3dhcmVob3VzZS5hdS5hdXRoMC5jb20vdXNlcmluZm8iXSwiaWF0IjoxNjY3Njg4Mzk2LCJleHAiOjE2Njc3NzQ3OTYsImF6cCI6ImlyRUkxZFZKT2J2aUo2N21ac1lGRWpQOTNFSWx0aHlzIiwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCB3cml0ZTpjYXRlZ29yeSByZWFkOmNhdGVnb3J5IHdyaXRlOml0ZW0gcmVhZDppdGVtIG9mZmxpbmVfYWNjZXNzIn0.kjogNt4dMxM624tnl8lBUXQEVxa6DirRiaBpVForUt7NcDttCPaiiZPEm_uo0QEvNy3ERrzMHZYwYmrXluiDSoelqynDbe88w5PpLZeFHPn-CA0gRzv0Y2zcqI-De19fFyjcrFwV8rSgH_MOq4rXMgOKTkDmIe8NUl1gSL98771nQ8b6--YHLKSBoeNuqTFBrUamQ3mfggUouK-kd75tnW754sORnwpS9_ujthim_Fnjfd9v-RtQBP3SemkpkuJdNKpUas_k7cs7HkMGQIiT42wpa8DVtggUh0NotRdmnLca1Qb4yT4PGwnJUi_XnUd-l8CXTer8HT5UWk8HEr-lHg',
      'Content-Type': 'application/json',
    },
  };

  http.post(url, payload, params);
}
