using UnityEngine;
using System.Collections;
//using Google.Protobuf.Gt;
using Google.Protobuf;
using System.Net.Sockets;
using System.Threading;
using System;
using System.Collections.Generic;
using System.Net;
public class GameTcpClient : MonoBehaviour {


	enum SocketState{
		None,           //无
		Connecting,     //链接中
		ConnectSuccess, //链接成功
		ConnectFail,    //链接失败
		ConnectEnd,     //链接完成
		SocketError,    //出现错误
		CloseConnect,   //关闭链接
	}
	private const int PACKET_HEAD_SIZE =2; //包头 2 ＋ 2 ＋ 2 (len, id, msg)
	//网络状态
	private SocketState m_socketState;
	//
	private Socket m_socket;
	//是否连接中
	public bool m_isConnect = false;
	//是否需要关闭连接
	private bool m_needClose = false;
	//
	private Action<bool> m_callbackOnConnect;
	//
	private Action m_callbackOnDisConnect;
	//发包
	private System.Object m_sendObject;
	private ManualResetEvent m_sendEvent;
	private LinkedList<IMessage> m_SendPackList;

	//收包
	private System.Object m_receiveObject;
	private LinkedList<IMessage> m_ReceivePackList;

	private Action<bool> callbackOnConnect;
	void Start () {
	
	}
	
	// Update is called once per frame
	void Update () {
		OnRun ();
	}
	public void OnRun(){

		if (m_socketState == SocketState.ConnectSuccess) {
			
			if (m_callbackOnConnect != null) {

				m_callbackOnConnect (true);
			}
	
			m_socketState = SocketState.ConnectEnd;

		} else if (m_socketState == SocketState.ConnectFail) {
			
			if (m_callbackOnConnect != null) {

				m_callbackOnConnect (false);
			}
			m_socketState = SocketState.ConnectEnd;

		} else if (m_socketState == SocketState.SocketError) {
			
			if (m_callbackOnDisConnect != null) {
				
				m_callbackOnDisConnect ();
			}
			m_socketState = SocketState.ConnectEnd;
		}
		if (m_isConnect == false) {

			return;
		}
		if (m_needClose == true) {

			Close ();

			return;
		}


		//发包队列
		if (m_ReceivePackList != null && m_ReceivePackList.Count != 0) {
			
			Queue<IMessage> sendQueue = new Queue<IMessage> ();

			lock (m_receiveObject) {
				
				while (m_ReceivePackList.Count > 0) {
					
					sendQueue.Enqueue (m_ReceivePackList.First.Value);

					m_ReceivePackList.RemoveFirst ();
				}
			}
			while(sendQueue.Count > 0){

				GameMessagehandler.MessageDispatch (sendQueue.Dequeue ());
			}
		}
	}

	//初始化网络
	public bool init(string ip, string port, Action<bool> callback){
	
		if (m_isConnect) {
			
			Close ();
		}

		try{

			Debug.Log("连接服务器 " + ip + " " + port);

			callbackOnConnect = callback;

			IPEndPoint ipe = new IPEndPoint (IPAddress.Parse (ip), Convert.ToInt32 (port));

			m_socket = new Socket (AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);

			m_socket.BeginConnect (ipe, new AsyncCallback (OnSocketConnectResult), m_socket);

			m_socketState = SocketState.Connecting;
		}
		catch(Exception e){
			
			Debug.LogError ("连接服务器失败!!!"+e);

			m_isConnect = false;

			m_socketState = SocketState.ConnectFail;

			return false;
		}

		return true;
	}

	//网络连接结果
	private void OnSocketConnectResult(IAsyncResult result){
		
		Socket s = (Socket)result.AsyncState;

		if (s.Connected) {
			
			m_isConnect = true;

			InitNetLoopThread();

			m_socketState = SocketState.ConnectSuccess;

			callbackOnConnect (true);

		} else {

			m_socketState = SocketState.ConnectFail;

			callbackOnConnect (false);
		}
	}

	//初始化循环线程
	private void InitNetLoopThread(){
		Debug.Log ("[InitNetLoopThread]");

		m_SendPackList = new LinkedList<IMessage> ();

		m_sendObject = new object ();

		m_sendEvent = new ManualResetEvent (false);

		m_ReceivePackList = new LinkedList<IMessage> ();

		m_receiveObject = new object ();

		Thread receiveThread = new Thread (ReceiveLoopThread);

		receiveThread.Start ();

		Thread sendThread = new Thread (SendLoopThread);

		sendThread.Start ();

	}
	//暂存容器
	byte[] piece;
	//消息长度
	int msglen = 0;
	//消息id

	//收包轮循
	private void ReceiveLoopThread(){
		
		Debug.Log("[ReceiveLoopThread]");

		while (m_socket != null && m_isConnect == true &&  m_needClose == false && m_socket.Connected) {
			try{
				if (m_socket.Available > 0) {

					if (piece == null && msglen == 0 && m_socket.Available > PACKET_HEAD_SIZE) { //接收包头
						
						piece = new byte[PACKET_HEAD_SIZE];
						//获取头 2 个字节
						m_socket.Receive (piece, PACKET_HEAD_SIZE, SocketFlags.None);

						Array.Reverse(piece);
						msglen = BitConverter.ToInt16(piece, 0);

					}
					if (piece != null && msglen > 0 ) { //接收包
						
						byte[] bodyBytes = new byte[msglen];

						int receiveBuffSize = 0;

						while(receiveBuffSize < msglen && m_isConnect == true){
							
							int tempSize = m_socket.Receive (bodyBytes,receiveBuffSize, msglen - receiveBuffSize, SocketFlags.None);

							receiveBuffSize += tempSize;

							Thread.Sleep(1);
						}
						if(receiveBuffSize == msglen){

							var idBytes = new byte[2]{bodyBytes[0], bodyBytes[1]};
							Array.Reverse(idBytes);
							UInt16 msgId =(UInt16) BitConverter.ToInt16(idBytes,0);

							var dataBytes = new byte[bodyBytes.Length - 2];
							System.Buffer.BlockCopy(bodyBytes, 2, dataBytes, 0, dataBytes.Length);
							IMessage packet = GameMessagehandler.DeserializePacket(dataBytes,msgId);

							lock (m_receiveObject) {
								
								m_ReceivePackList.AddLast (packet);
							}

							piece = null;

						}else{
							
							Debug.LogError ("ReceiveLoopThread error 接收包长度错误");

							m_needClose = true;
						}
					}
				} else {
					
					Thread.Sleep (10);
				}
			}
			catch(Exception e){
				
				Debug.LogError ("ReceiveLoopThread error "+e);

				m_needClose = true;
			}
		}
		Debug.Log ("ReceiveLoopThread die");
	}

	//发包轮循
	private void SendLoopThread(){
		
		Debug.Log("[SendLoopThread]");

		while (m_isConnect) {
			
			Queue<IMessage> sendQueue = new Queue<IMessage> ();

			lock (m_sendObject) {
				
				while (m_SendPackList.Count > 0) {
					
					sendQueue.Enqueue (m_SendPackList.First.Value);

					m_SendPackList.RemoveFirst ();
				}
			}
			while (sendQueue.Count > 0) {
				
				IMessage packet = null;

				packet = sendQueue.Dequeue();

				if (packet != null) {
					
					try{

						int msgId = GameMessagehandler.RequsetIdsByName[packet.Descriptor.Name];

						var loads = packet.ToByteArray ();

						var lenBytes = BitConverter.GetBytes((ushort)( 2 + loads.Length )); 
						var idBytes = BitConverter.GetBytes((ushort)msgId); 
						Array.Reverse(lenBytes);
						Array.Reverse(idBytes);
						var lenght = lenBytes.Length + idBytes.Length + loads.Length;

						var result = new byte[lenght];

						System.Buffer.BlockCopy(lenBytes, 0, result, 0, lenBytes.Length);
						System.Buffer.BlockCopy(idBytes, 0, result, lenBytes.Length, idBytes.Length);
						System.Buffer.BlockCopy(loads, 0, result, lenBytes.Length+idBytes.Length, loads.Length);

						m_socket.Send(result,0,lenght,SocketFlags.None);
					}
					catch(Exception e){
						
						Debug.Log ("Send error.. "+e.ToString());
					}
				}
			}
			m_sendEvent.Reset ();

			m_sendEvent.WaitOne ();
		}
	}

	public void SendPacket(IMessage packet){
		
		lock (m_sendObject) {
			
			m_SendPackList.AddLast (packet);
		}
		m_sendEvent.Set ();
	}

	public void Close(){
		
		m_isConnect = false;

		if (m_socket != null) {

			m_socket.Close();

			m_socket = null;
		}

	}

	void OnDestroy()
	{
		Close ();
	}
}
