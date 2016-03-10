send_gcm:
	curl -d 'tokens=gcm:APA91bGDU9E5oBZLEeiSTqNURgTAoR9TogY3CVarQT6hfEJDkwoGTKrxmTOdYZmvnX8FJtVTXxdkF35vo2Gw62cbucLHqrkKUk5yD_T-qnEwJ0MRxE_75swnTTJWFVkYorap-DRCqCKnF-2tbJO_MOt527ICH88moQ&payload={"m":"hello","h":"header","t":"m"}' localhost:5601/send

send_apns:
	curl  -d 'tokens=11218e477740b562f7702faf50db259628848f8f823c0732bfd36a19aae67100&payload={"aps":{"alert":{"body":"Brandon arrived at Mckesson"},"content_available":1,"custom":{"t":"gfi","u":"3be48e4dc9e7445fbaede8a28351d8b8:20c8d7e506214c6aa35586b68147ad83","g":"3b55767d-5e2d-cd83-84a3-b2cf64b36874","n":"gfi"},"sound":"p.caf","badge":13,"category":"SMH"}}' localhost:5601/send
