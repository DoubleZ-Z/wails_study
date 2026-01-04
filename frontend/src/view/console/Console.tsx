import './Console.less'
import React from "react";

const WorksheetManagement: React.FC = () => {
    return (
        <div className="console">
            <div className="console-title">
                console
            </div>
            <div className="console-content">
                <div className="console-item">
                    <div className="console-item-title">
                        <div className="console-item-title-text">
                            Worksheet Management
                        </div>
                    </div>
                    <div className="console-item-content">
                        <div className="console-item-content-text">
                            <div className="console-item-content-text-item"></div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default WorksheetManagement;