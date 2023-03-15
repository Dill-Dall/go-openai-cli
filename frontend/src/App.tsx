import React, { useState, useEffect } from "react";
import ReactMarkdown from "react-markdown";
import { MoonLoader } from "react-spinners";
import TextareaAutosize from 'react-textarea-autosize';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPaperPlane } from '@fortawesome/free-solid-svg-icons';

import "./App.css";
import "./global.css";

interface Conversation {
  id: string;
  name: string;
  messages: string[];
}

const App: React.FC = () => {
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [conversations, setConversations] = useState<Conversation[]>([
    { id: "conversation_0", name: "conversation", messages: [] },
  ]);
  const [selectedConversationIndex, setSelectedConversationIndex] = useState(0);
  const [selectedConversationMessages, setSelectedConversationMessages] = useState<string[]>([]);
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const newWs = new WebSocket("ws://localhost:8080/ws");
    setWs(newWs);
    return () => newWs.close();
  }, []);

  useEffect(() => {
    setSelectedConversationMessages(conversations[selectedConversationIndex].messages);
  }, [selectedConversationIndex, conversations]);

  useEffect(() => {
    // Create a WebSocket connection for fetching conversations
    const conversationsWs = new WebSocket("ws://localhost:8080/ws/conversations");

    // Set up the onmessage event to receive conversations
    conversationsWs.onmessage = (event: MessageEvent) => {
      const fetchedConversations = JSON.parse(event.data);
      setConversations(fetchedConversations);
    };

    // Close the WebSocket connection when the component unmounts
    return () => conversationsWs.close();
  }, []);

  useEffect(() => {
    if (!ws) return;

    ws.onmessage = (event: MessageEvent) => {
      setLoading(false);
      setConversations((prevConversations) => {
        return prevConversations.map((conversation, i) => {
          if (i === selectedConversationIndex) {
            return {
              ...conversation,
              messages: [...conversation.messages, event.data],
            };
          } else {
            return conversation;
          }
        });
      });
    };
  }, [ws, selectedConversationIndex]);


  const handleConversationClick = (index: number) => {
    setSelectedConversationIndex(index);
  };

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    if (!input || !ws) return;

    // Add input as a message to the current conversation
    setConversations((prevConversations) => {
      return prevConversations.map((conversation, i) => {
        if (i === selectedConversationIndex) {
          return {
            ...conversation,
            messages: [...conversation.messages, input],
          };
        } else {
          return conversation;
        }
      });
    });

    ws.send(input);
    setLoading(true);
    setInput("");
  };


  const handleAddConversation = () => {
    const newConversation: Conversation = {
      id: `conversation_${conversations.length}`,
      name: `conversation ${conversations.length}`,
      messages: [],
    };
    setConversations((prevConversations) => [...prevConversations, newConversation]);
    setSelectedConversationIndex(conversations.length);
  };

  return (
    <div className="App">
      <div className="sidebar">
        <button className="add-conversation-button" onClick={handleAddConversation}>
          New conversation
        </button>
        {conversations.map((conversation, i) => (
          <div
            key={conversation.id}
            className={`conversation ${i === selectedConversationIndex ? "active" : ""
              }`}
            onClick={() => handleConversationClick(i)}
          >
            {conversation.name}

          </div>
        ))}

      </div>
      <div className="main">
        <div className="messages">
          {selectedConversationMessages.map((message, i) => (
            <div
              key={i}
              className={`message ${i % 2 === 0 ? "even-message" : "odd-message"}`}
            >
              <ReactMarkdown>{message}</ReactMarkdown>
            </div>
          ))}
        </div>
        <form onSubmit={handleSubmit}>
          <div className="input-container">
            <TextareaAutosize
              value={input}
              onChange={(event) => setInput(event.target.value)}
              minRows={1}
              maxRows={6}
            />
            <button type="submit" disabled={loading} className="send-button">
              {loading ? (
                <MoonLoader size={15} color="#70b4e8" />
              ) : (
                <FontAwesomeIcon icon={faPaperPlane} />
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default App;
