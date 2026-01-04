import {useState} from 'react';
import './App.css';
import {Greet} from "../wailsjs/go/main/App";
import Console from "./view/console/Console";

function App() {
    const [resultText, setResultText] = useState("Please enter your name below 👇");
    const [name, setName] = useState('');
    const updateName = (e: any) => setName(e.target.value);
    const updateResultText = (result: string) => setResultText(result);

    function greet() {
        Greet(name).then(updateResultText);
    }

    return (
        <Console/>
    )
}

export default App
