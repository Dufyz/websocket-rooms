import { useState } from "react";
import useSocket from "./hooks/useSocket.hook";

export default function App() {
  const { connected, currentRooms, joinRoom, leaveRoom, sendMessage } =
    useSocket();

  const [message, setMessage] = useState<string>("");
  const [roomIdToConnect, setRoomIdToConnect] = useState<string>("");
  const [roomIdToMessage, setRoomIdToMessage] = useState<string>("");

  return (
    <div className="container">
      <h1>WebSocket Rooms</h1>

      <div className="chat-container">
        <div className="connection-area">
          <h2>Conectar à Sala</h2>
          <div className="input-group">
            <input
              type="text"
              value={roomIdToConnect}
              onChange={(e) => setRoomIdToConnect(e.target.value)}
              placeholder="ID da sala (ex: chat_1)"
              disabled={!connected}
            />
          </div>
          <div className="button-group">
            <button
              className="join-button"
              onClick={() => joinRoom(roomIdToConnect)}
              disabled={!connected || !roomIdToConnect}
            >
              Entrar na Sala
            </button>
            <button
              className="leave-button"
              onClick={() => leaveRoom(roomIdToConnect)}
              disabled={!connected || !roomIdToConnect}
            >
              Sair da Sala
            </button>
          </div>
        </div>

        <div className="message-area">
          <h2>Enviar Mensagem</h2>
          <div className="input-group">
            <select
              value={roomIdToMessage}
              onChange={(e) => setRoomIdToMessage(e.target.value)}
              disabled={!connected || currentRooms.size === 0}
            >
              <option value="" disabled>
                Selecione uma sala
              </option>
              {Array.from(currentRooms).map((room) => (
                <option key={room} value={room}>
                  {room}
                </option>
              ))}
            </select>
          </div>
          <div className="message-input-group">
            <input
              type="text"
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              placeholder="Digite sua mensagem"
              disabled={!connected || !roomIdToMessage}
            />
            <button
              className="send-button"
              onClick={() => {
                sendMessage(roomIdToMessage, message);
                setMessage("");
              }}
              disabled={!connected || !message || !roomIdToMessage}
            >
              Enviar
            </button>
          </div>
        </div>

        <div className="connection-status">
          <div
            className={`status-indicator ${
              connected ? "connected" : "disconnected"
            }`}
          />
          <span>{connected ? "Conectado" : "Desconectado"}</span>
        </div>

        {currentRooms.size > 0 && (
          <div className="rooms-list">
            <h2>Salas Conectadas:</h2>
            <ul>
              {Array.from(currentRooms).map((room) => (
                <li key={room}>{room}</li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </div>
  );
}
