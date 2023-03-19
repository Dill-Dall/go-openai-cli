import React, { useState, useEffect, ChangeEvent, useRef } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faTrashAlt } from "@fortawesome/free-solid-svg-icons";
import { config, mockConversations, mockSystemodels } from "./config";

import "./App.css";
import "./global.css";
import Main from "./components/Main";
import SystemModelSelect from "./components/SystemModelSelect";

interface Conversation {
  id: string;
  name: string;
  messages: string[];
}


const useChat = () => {
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  return { messagesEndRef, textareaRef };
};


const App: React.FC = () => {
  const [conversations, setConversations] = useState<Conversation[]>((config.useMockData ? mockConversations : []));
  const [selectedConversationIndex, setSelectedConversationIndex] = useState(-1);
  const [selectedConversationMessages, setSelectedConversationMessages] = useState<string[]>([]);
  const [shouldCreateNewConversation, setShouldCreateNewConversation] = useState(true);
  const [isAddingConversation, setIsAddingConversation] = useState(false);


  const [selectedSystemModel, setSelectedSystemModel] = useState<string>(
    config.useMockData ? mockSystemodels[0] : ""
  );

  const { messagesEndRef, textareaRef } = useChat();

  const [input, setInput] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    fetch("http://localhost:8080/api/conversations")
      .then((res) => res.json())
      .then((data) => setConversations(data));
  }, []);

  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, [messagesEndRef, selectedConversationMessages]);


  const handleConversationClick = (index: number) => {
    setSelectedConversationIndex(index);
    setSelectedConversationMessages(conversations[index]?.messages || []);
    setShouldCreateNewConversation(false);
  };

  const handleAddConversation = () => {
    setIsAddingConversation(true);
    setSelectedConversationIndex(-1);
    setSelectedConversationMessages([]);
    setShouldCreateNewConversation(true);
  };


  const handleDelete = (id: string) => {
    fetch("http://localhost:8080/api/conversations", {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(id),
    }).then(() => {
      const updatedConversations = conversations.filter((conversation) => conversation.id !== id);
      setConversations(updatedConversations);
      handleAddConversation()
    })
      .catch((error) => {
        console.error("Error deleting conversation:", error);
      });
  };

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    if (!input) {
      if (isAddingConversation) {
        setIsAddingConversation(false);
      }
      if (textareaRef.current) {
        textareaRef.current.focus();
      }
      return;
    }
    let id = "";
    if (selectedConversationIndex >= 0) {
      id = conversations[selectedConversationIndex].id;
    }
    const messageData = {
      content: input,
      id: id,
      system_model: selectedSystemModel,
    };

    setInput("")
    setIsLoading(true);

    fetch("http://localhost:8080/api/conversations", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(messageData),
    })
      .then((res) => res.json())
      .then((updatedConversation) => {
        console.log("Received updated conversation:", updatedConversation);
        if (shouldCreateNewConversation) {
          const updatedConversations = [...conversations, updatedConversation];
          console.log("Updated conversations (new):", updatedConversations);
          setConversations(updatedConversations);
          setSelectedConversationIndex(updatedConversations.length - 1);
          setShouldCreateNewConversation(false);
        } else {
          const updatedConversations = conversations.map((conversation, i) => {
            if (i === selectedConversationIndex) {
              return updatedConversation;
            } else {
              return conversation;
            }
          });
          console.log("Updated conversations (existing):", updatedConversations);
          setConversations(updatedConversations);
        }
        setSelectedConversationMessages(updatedConversation.messages);
      })

      .finally(() => {
        setIsLoading(false);
      });
  }

  return (
    <div className="App">
      <div className="sidebar">
        <div
          className={`conversation ${selectedConversationIndex === -1 ? "active" : ""
            }`}
          onClick={handleAddConversation}
        >
          New conversation
        </div>
        {conversations.map((conversation, i) => (
          <ConversationItem
            key={conversation.id}
            conversation={conversation}
            isSelected={i === selectedConversationIndex}
            onClick={() => handleConversationClick(i)}
            onDelete={() => handleDelete(conversation.id)}
          />
        ))}
      </div>
      <Main
        selectedConversationMessages={selectedConversationMessages}
        messagesEndRef={messagesEndRef}
        handleSubmit={handleSubmit}
        textareaRef={textareaRef}
        input={input}
        setInput={setInput}
        isLoading={isLoading} />
      <SystemModelSelect
        selectedSystemModel={selectedSystemModel}
        setSelectedSystemModel={setSelectedSystemModel} />
    </div >
  );
};

export default App;



interface ConversationItemProps {
  conversation: Conversation;
  isSelected: boolean;
  onClick: () => void;
  onDelete: () => void;
}

const ConversationItem: React.FC<ConversationItemProps> = ({
  conversation,
  isSelected,
  onClick,
  onDelete,
}) => {
  return (
    <div
      className={`conversation ${isSelected ? "active" : ""}`}
      onClick={onClick}
    >
      <div>{conversation.name}</div>
      <div className="delete-button" onClick={onDelete}>
        <FontAwesomeIcon icon={faTrashAlt} />
      </div>
    </div>
  );
};