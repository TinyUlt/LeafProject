using UnityEngine;
using System.Collections;
using System;
using MiniJSON;
using System.Collections.Generic;
public class GameHttpClient : MonoBehaviour {

	public NetInfo gameInfo;

	public delegate void Action<T1, T2, T3, T4, T5> (T1 arg1, T2 arg2, T3 arg3, T4 arg4, T5 arg5);

	// Use this for initialization
	void Start () {
	
	}
	
	// Update is called once per frame
	void Update () {
	
	}
	public void NetConnect(Action<bool, int, string, string, string> callback){

		Debug.Log ("[NetConnect]");

		//测试本机网络
		StartCoroutine (CheckNetConnect((netEnable) => {
			
			if(netEnable){
				//链接服务器
				StartCoroutine (ServerNetConnect((isConnect, requrestData) => {

					string sessionid = "";

					string uid = "";

					string gateway = "";

					int code = -1;

					if(isConnect){
						//解析json
						var onePathDic = Json.Deserialize (requrestData) as Dictionary<string, object>;

						code = Convert.ToInt32( onePathDic ["Code"] );

						if(code == 0){

							var User = Convert.ToString( onePathDic ["User"] );
							var Psw = Convert.ToString( onePathDic ["Psw"] );

							Debug.Log (User);
							Debug.Log (Psw);
//							var data = onePathDic ["data"] as Dictionary<string, object>;
//
//							gateway = Convert.ToString (data ["gateway"]);
//
//							uid = Convert.ToString(data ["uid"]);
//
//							sessionid = Convert.ToString(data ["sessionid"]);
						}
					}
					//返回给登录脚本
					callback(isConnect, code, sessionid, uid, gateway);
				}));
			}
		}));
	}
	//测试本机网络
	public IEnumerator CheckNetConnect(Action<bool> callback){
		//http get 方法
		WWW www = new WWW (gameInfo.HttpCheck);

		yield return www;
		//检测是否连接成功
		if (!string.IsNullOrEmpty (www.error)) {
			
			callback (false);

		} else {
			
			callback (true);
		}
	}
	//连接服务器网络
	public IEnumerator ServerNetConnect(Action<bool, string> callback){

		WWWForm from = new WWWForm ();

//		string cmd = "1030";
//		string channel = "2";
//		string login_name = "whx";
//		string password = "e10adc3949ba59abbe56e057f20f883e";
//		string time = secondsSinceEpoch().ToString();
//		string sha1 =  "channel="+channel+"%cmd="+cmd+"%login_name="+login_name+"%password="+password+"%t="+time+gameInfo.HttpKey;
//		sha1 = SHA1(sha1).ToLower();


		from.AddField ("User", "wuhongxing");
		from.AddField ("Psw", "123454321");
//		from.AddField ("login_name", login_name);
//		from.AddField ("password", password);
//		from.AddField ("t", time);
		//http post 方法
		WWW www = new WWW (gameInfo.HttpHost+"/login",from);

		yield return www;
		//检测是否连接成功
		if (!string.IsNullOrEmpty (www.error)) {

			callback (false, "");

		} else {

			var returnData = System.Text.Encoding.UTF8.GetString (www.bytes);

			callback (true, returnData);
		}
	}
	//sha1签名
	public static string SHA1(string text)
	{
		byte[] cleanBytes = System.Text.Encoding.Default.GetBytes(text);

		byte[] hashedBytes = System.Security.Cryptography.SHA1.Create().ComputeHash(cleanBytes);

		return System.BitConverter.ToString(hashedBytes).Replace("-", "");
	}
	//获取时间
	public static int secondsSinceEpoch(){

		System.TimeSpan t = System.DateTime.UtcNow - new System.DateTime(1970, 1, 1);

		return (int)t.TotalSeconds;
	} 
}
