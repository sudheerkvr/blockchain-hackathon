console.log("Started");
var obj, dbParam, xmlhttp, myObj, x, message, txt = "";
var chainHeight;
var penaltyAmt = 0;
obj = {
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "dda20452e357d3cedaa6dbf430653d16388215449935b018f2b102e66f51843c83e5650eb19d5297ba9cba13a8276aeead243dd85155efcb4c77ec24ab386f74"
    },
    "ctorMsg": {
      "function": "get_pending_amount",
      "args": [
        "Primetime Editing Services"
      ]
    },
    "secureContext": "test_user0"
  },
  "id": 2
};
dbParam = JSON.stringify(obj);
console.log(dbParam);
xmlhttp = new XMLHttpRequest();
xmlhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
		console.log("this.responseText : "+this.responseText);
        returnObj = JSON.parse(this.responseText);
		console.log("returnObj1 : "+returnObj.result.message);
		message = JSON.parse(returnObj.result.message);
		console.log("message : "+message.completedworkamount);
		workTrans = message.projectresults;
		paidTrans = message.amounttransactions;
		allTrans = [];
		for(x in workTrans)
		{
			allTrans.push({"name" : workTrans[x].name,
							"date" : workTrans[x].date,
							"amount" : workTrans[x].derivedamount,
							"projectname" : workTrans[x].projectname,
							"taskname" : workTrans[x].taskname,
							"type" : "WORK_DONE"
			});
		}
		for(x in paidTrans)
		{
			allTrans.push({"name" : "Amount Received",
							"date" : paidTrans[x].date,
							"amount" : paidTrans[x].amountpaid,
							"projectname" : paidTrans[x].projectname,
							"type" : "AMOUNT_PAID"
			});
		}
		console.log(JSON.stringify(allTrans));
		
		for(x in allTrans) {
			var d = allTrans[x].date.split("/");
			console.log(d[0] + " : "+d[1]+" : "+d[2]);
			year = Number(d[2]);
			month = Number(d[0]) - 1;
			day = Number(d[1]);
			allTrans[x]["dateObj"] = new Date(year, month, day);
		}
		allTrans.sort(function(a,b) {
			return a.dateObj.valueOf() < b.dateObj.valueOf();
		});
		console.log(JSON.stringify(allTrans));
        txt += "<table class='mytable mytable-hover'><tbody>"
        for (x in allTrans) {
			if(allTrans[x].type == "AMOUNT_PAID") {
				txt += "<tr class='amountPaidRow'>";
			} else {
				txt += "<tr>";
			}
            txt += "<td class='trxData'><b><span class='trxDesc'>" + allTrans[x].name + "</span></b><br/>" + allTrans[x].projectname;
			if((allTrans[x].taskname != null) && (allTrans[x].taskname != "")) {
				txt += " > " + allTrans[x].taskname;
			}
			txt += "</td><td class='trxData'>"
				+ allTrans[x].date + "</td><td class='trxData'><b><span class='myAmountOne'>"
				+ formatAmountZero(parseInt(allTrans[x].amount)) + "</span>";
			if((allTrans[x].projectname == "Making of Big Labowski Project") && (allTrans[x].taskname == "")) {
				txt += "<br/><span style='color:rgb(243,0,0)'>Penalty: $1K</span></b></td>";
				penaltyAmt = 1000;
			} else {
				txt += "</b></td>";
			}
				
			txt +="<td class='trxData'><img src='img/logo.svg' alt='Valid' height='30' width='30'/></td></tr>";
        }
        txt += "</tbody></table>"
        document.getElementById("transactions").innerHTML = txt;

		var canvas=document.getElementById("contractAmountCanvas");
		var ctx=canvas.getContext("2d");
		var colors=['#0D4A92','#D8D8D8'];
		completedworkamount = parseInt(message.completedworkamount);
		pendingcontractamount = parseInt(message.pendingcontractamount);
		txt1 = "<table><tr><td><span class='mylabel'>Completed</span></td><td align='right'><span class='mylabel myAmountZero'>" + formatAmountZero(completedworkamount) + "</span></td></tr><tr><td><span class='mylabel'>Remaining</span></td><td align='right'><span class='mylabel myAmountZero'>" + formatAmountZero(pendingcontractamount) + "</span></td></tr></table>";
        document.getElementById("WorkData").innerHTML = txt1;
		var workCompletedPercent;
		if((completedworkamount + pendingcontractamount) == 0) {
			workCompletedPercent = 0;
		} else {
			workCompletedPercent = (completedworkamount * 100) / (completedworkamount + pendingcontractamount);
		}
		workCompletedPercent = Math.round(workCompletedPercent);
		var values=[completedworkamount,pendingcontractamount];
		console.log("values : "+values);
		var labels=['Completed Work Amount','Pending Contract Amount'];
		
		dmbChart(125,125,90,30,values,colors,labels,0, ctx, workCompletedPercent);

		paidAmountCanvas=document.getElementById("paidAmountCanvas");
		ctx=paidAmountCanvas.getContext("2d");
		colors=['#0D4A92','#D8D8D8'];
		amountpaid = parseInt(message.amountpaid);
		balancetobepaid = parseInt(message.balancetobepaid);
		balancetobepaid -= penaltyAmt;
		txt1 = "<table><tr><td><span class='mylabel'>Amount Paid</span></td><td align='right'><span class='mylabel myAmountZero'>" + formatAmountZero(amountpaid) + "</span></td></tr><tr><td><span class='mylabel'>Remaining Balance</span></td><td align='right'><span class='mylabel myAmountZero'>" + formatAmountZero(balancetobepaid) + "</span></td></tr></table>";
        document.getElementById("BalanceData").innerHTML = txt1;
		var amountPaidPercent;
		if((amountpaid + balancetobepaid) == 0) {
			amountPaidPercent = 0;
		} else {
			amountPaidPercent = (amountpaid * 100) / (amountpaid + balancetobepaid);
		}
		amountPaidPercent = Math.round(amountPaidPercent);
		values=[amountpaid,balancetobepaid];
		console.log("values : "+values);
		labels=['Amount paid','Balance Amount to be paid'];
		
		dmbChart(125,125,90,30,values,colors,labels,0, ctx, amountPaidPercent);
    }
};
xmlhttp.open("POST", "http://192.168.99.100:7050/chaincode", true);
xmlhttp.setRequestHeader("Content-type", "application/json");
xmlhttp.send(dbParam);
console.log("Request1 sent");

xmlhttp = new XMLHttpRequest();
xmlhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
		console.log("this.responseText2 : "+this.responseText);
        returnObj = JSON.parse(this.responseText);
		chainHeight = returnObj.height
		console.log("Chain Height : "+returnObj.height);
		getBlock(0);
    }
};
xmlhttp.open("GET", "http://192.168.99.100:7050/chain", true);
xmlhttp.setRequestHeader("Content-type", "application/json");
xmlhttp.send();
console.log("Request2 sent");

function dmbChart(cx,cy,radius,arcwidth,values,colors,labels,selectedValue, ctx, percentage){
    var tot=0;
    var accum=0;
    var PI=Math.PI;
    var PI2=PI*2;
    var offset=-PI/2;
    ctx.lineWidth=arcwidth;
    for(var i=0;i<values.length;i++) {
		tot+=values[i];console.log("tot : "+tot);
	}
	if(tot == 0) {
		values[values.length - 1] = 100;
		tot = 100;
		console.log("values.length : "+values.length);
		console.log("values[values.length - 1] : "+values[values.length - 1]);
	}
    for(var i=0;i<values.length;i++){
        ctx.beginPath();
        ctx.arc(cx,cy,radius,
            offset+PI2*(accum/tot),
            offset+PI2*((accum+values[i])/tot)
        );
        ctx.strokeStyle=colors[i];
        ctx.stroke();
        accum+=values[i];
		var font = "50px verdana";
		ctx.font = font;
		var width = ctx.measureText(percentage+"%").width;
		var height = ctx.measureText("w").width;
		ctx.fillText(percentage+"%", 125 - (width/2) ,125 + (height/2));
    }
}

function formatAmount(amount) {
	amount = Math.round(amount / 100) / 10;
	str = amount.toString();
	if(str.indexOf(".") == -1) {
		str = str.concat(".0");
	}
	return "$"+ str + "K";
}

function formatAmountZero(amount) {
	amount = Math.round(amount/1000);
	return "$"+ amount + "K";
}

function getBlock(i){
	xmlhttp = new XMLHttpRequest();
	xmlhttp.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			console.log("this.responseText3 : "+this.responseText);
			returnObj = JSON.parse(this.responseText);
			var temp =  {
				id: i,
				blockstats: returnObj
			};
			console.log("Block Details :"+JSON.stringify(temp));
			new_block(temp);
			if ((i+1) < chainHeight) {
				getBlock(i+1);
			}
		}
	};
	xmlhttp.open("GET", "http://192.168.99.100:7050/chain/blocks/"+i, true);
	xmlhttp.setRequestHeader("Content-type", "application/json");
	xmlhttp.send();
	console.log("Request3 sent");
}

function onMessage(msg){
		try{
			var msgObj = JSON.parse(msg.data);
			if(msgObj.marble){
				console.log('rec', msgObj.msg, msgObj);
				build_ball(msgObj.marble);
			}
			else if(msgObj.msg === 'chainstats'){
				console.log('rec', msgObj.msg, ': ledger blockheight', msgObj.chainstats.height, 'block', msgObj.blockstats.height);
				if(msgObj.blockstats && msgObj.blockstats.transactions) {
                    var e = formatDate(msgObj.blockstats.transactions[0].timestamp.seconds * 1000, '%M/%d/%Y &nbsp;%I:%m%P');
                    $('#blockdate').html('<span style="color:#fff">TIME</span>&nbsp;&nbsp;' + e + ' UTC');
                    var temp =  {
                        id: msgObj.blockstats.height,
                        blockstats: msgObj.blockstats
                    };
                    new_block(temp);								//send to blockchain.js
				}
			}
			else console.log('rec', msgObj.msg, msgObj);
		}
		catch(e){
			console.log('ERROR', e);
		}
	}
