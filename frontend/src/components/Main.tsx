import { faPaperPlane } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { useEffect } from "react";
import ReactMarkdown from "react-markdown";
import { MoonLoader } from "react-spinners";
import TextareaAutosize from "react-textarea-autosize";
import "./main.css";

type MainProps = {
	selectedConversationMessages: string[];
	messagesEndRef: React.RefObject<HTMLDivElement>;
	textareaRef: React.RefObject<HTMLTextAreaElement>;
	handleSubmit: (event: React.FormEvent<HTMLFormElement>) => void;
	input: string;
	setInput: (input: string) => void;
	isLoading: boolean;
};

const Main: React.FC<MainProps> = ({
	selectedConversationMessages,
	messagesEndRef,
	textareaRef,
	handleSubmit,
	input,
	setInput,
	isLoading,

}) => {


	return (
		<div className="main">
			<div className="messages">
				{selectedConversationMessages.map((message, i) => (
					<div
						key={i}
						className={`message ${i % 2 === 0 ? "even-message" : "odd-message"}`}
					>
						<ReactMarkdown>{addaptedMessage(message)}</ReactMarkdown>
					</div>
				))}
				<div ref={messagesEndRef} />
			</div>
			<form onSubmit={handleSubmit}>
				<div className="input-container">
					<TextareaAutosize
						ref={textareaRef}
						value={input}
						onChange={(event) => setInput(event.target.value)}
						onKeyDown={(event) => handleTextareaKeyDown(event, handleSubmit)}
					/>


					<button type="submit" disabled={isLoading} className="send-button">
						{isLoading ? (
							<MoonLoader size={15} color="#70b4e8" />
						) : (
							<FontAwesomeIcon icon={faPaperPlane} />
						)}
					</button>
				</div>
			</form>
		</div>
	);
};



const handleTextareaKeyDown = (event: React.KeyboardEvent<HTMLTextAreaElement>, handleSubmit: (event: React.FormEvent<HTMLFormElement>) => void) => {
	if (event.key === "Enter" && !event.shiftKey) {
		event.preventDefault();
		handleSubmit(event as unknown as React.FormEvent<HTMLFormElement>);
	}
};




const addaptedMessage = (message: string) => {

	message = message.replace(/^(USER|ASSISTANT):\s*/, "");
	message = message.replace(/^Title:.*/, "");

	const urlRegex = /URL_\[([^\]]+)\]/g;
	message = message.replace(urlRegex, '![alt text]($1)');

	return message;
}

export default Main;