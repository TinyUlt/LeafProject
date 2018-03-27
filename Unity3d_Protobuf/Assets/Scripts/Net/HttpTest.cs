using UnityEngine;
using System.Collections;
using Google.Protobuf;
public class HttpTest : MonoBehaviour {
	
	public GameHttpClient httpClient;
	void Start() {
		
	}

	public void TryLogin(){

		Debug.Log ("[TryLogin]");

		httpClient.NetConnect ((isConnectNet, errorCode, sessionid, uid, gateway)=>{

			if(isConnectNet){

				Debug.Log("网络活跃");	



			}else{

				Debug.Log("网络不活跃，请检测网络!!!");				
			}
		});
	}
}
