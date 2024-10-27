# go-echo-practice
Go Echo × SQLboilerでTODOリストを作ってみる

## 試した所管
- Ginと比べてコード量が多くなるかも
	- controllerで戻り値を必ずerror型を返す必要があるため、きちんとレスポンスをreturnする必要がある
	- JSON型で返すときにはgin.Hのようなものがないため、自分で加工する必要がある
	- sqlboilerを使っていると、contextがecho.contextだと不適のため、別途でcontextを宣言する必要がある
- Ginと比べると拡張性が高いと言われていて規模が大きくなってくると役立つらしいが、いかに？
