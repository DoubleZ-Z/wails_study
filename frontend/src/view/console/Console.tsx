import './Console.less'
import React, {useState} from "react";
import {Button, Table} from "antd";
import {TcpMessageHistory} from "../../../wailsjs/go/main/App";
import {history} from "../../../wailsjs/go/models";
import {ColumnType} from "antd/es/table";
import {EventsOn} from "../../../wailsjs/runtime";
import MessageHistory = history.MessageHistory;


const Console: React.FC = () => {
    const [messageHistoryList, setMessageHistoryList] = useState<MessageHistory[]>([])


    EventsOn("messageHistoryReflush", function (res: MessageHistory[]) {
        console.log(res)
        setMessageHistoryList(res.sort((a, b) => b.traceId.localeCompare(a.traceId)))
    });

    function onclick() {
        TcpMessageHistory().then(res => {
            setMessageHistoryList(res.sort((a, b) => b.traceId.localeCompare(a.traceId)))
        })
    }

    const messageHistoryColumns = getMessageHistoryColumns();

    return (
        <div className="console">
            <div className="console-title">
                console
            </div>
            <Button onClick={onclick}>
                刷新
            </Button>
            <div className="table-container">
                <Table
                    size={"small"}
                    rowKey='traceId'
                    className={"history-table"}
                    columns={messageHistoryColumns}
                    dataSource={messageHistoryList}
                    pagination={false}
                    scroll={{x: 'max-content', y: 300}}
                />
            </div>
        </div>
    )
}

const getMessageHistoryColumns = (): ColumnType<MessageHistory>[] => {
    return [
        {
            title: 'traceId',
            dataIndex: 'traceId',
            key: 'traceId',
            align: 'center',
            width: 120,
            fixed: 'start',
        },
        {
            title: 'station',
            dataIndex: 'stationNo',
            key: 'stationNo',
            align: 'center',
            width: 80,
            fixed: 'start',
        },
        {
            title: 'category',
            dataIndex: 'category',
            key: 'category',
            align: 'center',
            width: 50,
        },
        {
            title: 'success',
            dataIndex: 'success',
            key: 'success',
            align: 'center',
            width: 50,
            render: (success: boolean) => {
                return success ? 'success' : 'fail'
            }
        },
        {
            title: 'action',
            dataIndex: 'action',
            key: 'action',
            align: 'center',
            width: 200,
        },
        {
            title: 'requestTime',
            dataIndex: 'requestTime',
            key: 'requestTime',
            align: 'center',
            width: 200,
        },
        {
            title: 'responseTime',
            dataIndex: 'responseTime',
            key: 'responseTime',
            align: 'center',
            width: 200,
        },
        {
            title: 'duration',
            dataIndex: 'duration',
            key: 'duration',
            align: 'center',
            width: 20,
        },
        {
            title: 'requestDelay',
            dataIndex: 'requestDelay',
            key: 'requestDelay',
            align: 'center',
            width: 20,
        }
    ]
}

export default Console;