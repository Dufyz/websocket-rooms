import { useEffect, useState } from "react";

type Message = {
  type: string;
  room_id: string;
  payload: any;
};

export default function useSocket() {
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [connected, setConnected] = useState<boolean>(false);
  const [currentRooms, setCurrentRooms] = useState<Set<string>>(new Set());

  const joinRoom = (room_id: string) => {
    if (!socket || socket.readyState !== WebSocket.OPEN) return;

    const event: Message = {
      type: "join",
      room_id,
      payload: {},
    };

    socket.send(JSON.stringify(event));

    setCurrentRooms((prev) => new Set(prev).add(room_id));

    console.log("Tentando entrar na sala:", event);
  };

  const leaveRoom = (room_id: string) => {
    if (!socket || socket.readyState !== WebSocket.OPEN) return;

    const event: Message = {
      type: "leave",
      room_id,
      payload: {},
    };

    socket.send(JSON.stringify(event));

    setCurrentRooms((prev) => {
      const newRooms = new Set(prev);
      newRooms.delete(room_id);
      return newRooms;
    });

    console.log("Saindo da sala:", event);
  };

  const sendMessage = (room_id: string, message: string) => {
    if (!socket || socket.readyState !== WebSocket.OPEN) return;

    const event = {
      type: "message",
      room_id,
      payload: {
        message,
      },
    };
    socket.send(JSON.stringify(event));

    console.log("Mensagem enviada:", event);
  };

  useEffect(() => {
    const ws = new WebSocket("ws://localhost:3000/api/web-socket");

    ws.onopen = () => {
      setConnected(true);
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        console.log("Mensagem recebida:", data);

        if (data.type === "roomJoined") {
          console.log("Entrou na sala:", data.payload);
        }
      } catch (error: any) {
        console.log("Mensagem recebida (não-JSON):", event.data);
      }
    };

    ws.onclose = () => {
      setConnected(false);
      setCurrentRooms(new Set());
    };

    setSocket(ws);

    return () => {
      ws.close();
    };
  }, []);

  return {
    socket,
    connected,
    currentRooms,
    joinRoom,
    leaveRoom,
    sendMessage,
  };
}
