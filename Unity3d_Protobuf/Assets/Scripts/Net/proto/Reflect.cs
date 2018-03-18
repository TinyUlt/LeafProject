using Msg;
using System;
//消息的定义
public class Reflect  {

	public static void Init(){
	
		GameMessagehandler.setReflect (typeof(Hello),  (UInt16)MessageId.Hello );
		GameMessagehandler.setReflect (typeof(UserLoginRequest),  (UInt16)MessageId.UserLoginRequest );
	
	}
}
