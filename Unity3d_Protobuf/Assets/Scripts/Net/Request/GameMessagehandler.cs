using UnityEngine;
using System.Collections;
using System.Collections.Generic;
using Msg;
using Google.Protobuf;
using System;
public static class GameMessagehandler {

	//发送
	public static Dictionary<Type, UInt16> RequsetIdsByType;

	public static Dictionary<string, UInt16> RequsetIdsByName;
	//接受
	public static Dictionary<UInt16, Type> MsgMaps;

	public static void init(){

		RequsetIdsByType = new Dictionary<Type, UInt16> ();

		RequsetIdsByName = new Dictionary<string, UInt16> ();

		MsgMaps = new Dictionary<UInt16, Type> ();

		Reflect.Init ();
	}

	public static void setReflect(Type type, UInt16 id ){

		RequsetIdsByType.Add (type, id);

		RequsetIdsByName.Add (type.Name, id);

		MsgMaps.Add (id, type);
	}

	public static IMessage DeserializePacket(byte[] data, UInt16 id){

		var type = MsgMaps [id];

		IMessage msg = System.Activator.CreateInstance (type) as IMessage;

		msg.MergeFrom (data);

		return msg;
	}

	//消息分发
	public static void MessageDispatch(IMessage msg){

		Debug.Log (msg.Descriptor.GetType());
		Debug.Log (msg.Descriptor.ContainingType);
		Debug.Log (msg.Descriptor.EnumTypes);
		Debug.Log (msg.Descriptor.NestedTypes);



		var descriptor = msg.Descriptor;
		foreach (var field in descriptor.Fields.InDeclarationOrder())
		{
			Debug.LogFormat(
				"Field {0} ({1}): {2}",
				field.FieldNumber,
				field.Name,
				field.Accessor.GetValue(msg));
		}

		var name = msg.Descriptor.Name;

		switch (name) {

		case "Hello":
			{
				Hello lr = msg as Hello;

				Debug.Log (lr);

				break;
			}
		case "UserLoginRequest":
			{
				UserLoginRequest lr = msg as UserLoginRequest;

				Debug.Log (lr);

				break;
			}
		default:
			{
				break;
			}
		}
	}
}
