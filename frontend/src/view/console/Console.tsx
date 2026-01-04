import './Console.less'
import React, {useState} from "react";
import {Button} from "antd";
import {TcpMessageHistory} from "../../../wailsjs/go/main/App";
import {Simulate} from "react-dom/test-utils";
import reset = Simulate.reset;

const WorksheetManagement: React.FC = () => {
    const [messageHistoryList, setMessageHistoryList] = useState<any>([])

    function onclick() {
        TcpMessageHistory().then(res=>{
            setMessageHistoryList(res)
        })
    }

    console.log(messageHistoryList)

    return (
        <div className="console">
            <div className="console-title">
                console
            </div>
            <Button onClick={onclick}>
                点击
            </Button>
            {messageHistoryList.map((item: any) => (
                <div>
                    {item.TraceId}
                </div>
            ))}
        </div>
    )
}

export default WorksheetManagement;