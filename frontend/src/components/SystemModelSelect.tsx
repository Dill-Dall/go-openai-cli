import { useEffect, useState } from "react";
import { config, mockSystemodels } from "../config";
import "./SystemModelSelect.css";

type SystemModelSelectProps = {
	selectedSystemModel: string;
	setSelectedSystemModel: (model: string) => void;
};

const SystemModelSelect: React.FC<SystemModelSelectProps> = ({
	selectedSystemModel,
	setSelectedSystemModel,
}) => {
	const [systemodels, setSystemodels] = useState<string[]>(
		config.useMockData ? mockSystemodels : []
	);


	useEffect(() => {
		fetch("http://localhost:8080/api/systemodels")
			.then((res) => res.json())
			.then((data) => {
				setSystemodels(data);
				setSelectedSystemModel(
					data.includes("AI") ? "AI" : data[0]
				);
			});
	}, []);


	return (
		<div className="system-models">
			<select
				value={selectedSystemModel}
				onChange={(event) => setSelectedSystemModel(event.target.value)}
				className="system-models-dropdown"
			>
				{systemodels.map((model, index) => (
					<option key={index} value={model}>
						{model}
					</option>
				))}
			</select>
		</div>
	);
};



export default SystemModelSelect;