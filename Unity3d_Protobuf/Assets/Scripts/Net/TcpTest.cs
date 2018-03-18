using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using Msg;
public class TcpTest : MonoBehaviour {

	public GameTcpClient tcpClient;
	// Use this for initialization
	void Start () {
		
	}
	
	// Update is called once per frame
	void Update () {
		
	}

	public void TryLogin(){

		Debug.Log ("[TryLogin]");

		bool isCreated = tcpClient.init("127.0.0.1", "3566", (isConnectTcp)=>{

			if(isConnectTcp){

				Debug.Log("tcp 服务器连接成功");	

				GameMessagehandler.init();

				UserLoginRequest request = new UserLoginRequest();

				request.UserName = "whx";

				tcpClient.SendPacket(request);

			}else{

				Debug.Log("tcp 服务器连接失败!!!");	
			}
		});
		if(isCreated){

			Debug.Log("创建网络成功");
		}else{

			Debug.Log("创建网络失败!!!");
		}
	}
}
