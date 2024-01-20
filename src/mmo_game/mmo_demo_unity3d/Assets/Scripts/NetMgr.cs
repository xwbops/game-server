using UnityEngine;
using System;
using System.Collections;
using System.Collections.Generic;

public class NetMgr : MonoBehaviour
{
    private SocketClient socket;

    private string url;

    //Init
    public static Action<int> OnBorn;
    //move
    public static Action<Pb.BroadCast> OnMove;
    //leave
    public static Action<int> OnOver;


    SocketClient SocketClient
    {
        get
        {
            if (socket == null)
                socket = new SocketClient();
            return socket;
        }
    }

    void Awake()
    {
        DontDestroyOnLoad(gameObject);
        Init();
    }

    void Init()
    {
        SocketClient.OnRegister();
    }


    void Update()
    {
    }
    public void SendConnect(string url, int port)
    {
        SocketClient.SendConnect(url, port);
    }


    public void SendMessage(ByteBuffer buffer)
    {
        SocketClient.SendMessage(buffer);
    }

    new void OnDestroy()
    {
        SocketClient.OnRemove();
    }
}
