export namespace history {
	
	export class MessageHistory {
	    traceId: string;
	    stationNo: string;
	    category?: string;
	    success?: boolean;
	    resultCode?: number;
	    resultMessage?: string;
	    action?: string;
	    requestMessage?: string;
	    responseMessage?: string;
	    requestTime?: string;
	    receiveTime?: string;
	    responseTime?: string;
	    duration?: number;
	    requestDelay?: number;
	
	    static createFrom(source: any = {}) {
	        return new MessageHistory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.traceId = source["traceId"];
	        this.stationNo = source["stationNo"];
	        this.category = source["category"];
	        this.success = source["success"];
	        this.resultCode = source["resultCode"];
	        this.resultMessage = source["resultMessage"];
	        this.action = source["action"];
	        this.requestMessage = source["requestMessage"];
	        this.responseMessage = source["responseMessage"];
	        this.requestTime = source["requestTime"];
	        this.receiveTime = source["receiveTime"];
	        this.responseTime = source["responseTime"];
	        this.duration = source["duration"];
	        this.requestDelay = source["requestDelay"];
	    }
	}

}

