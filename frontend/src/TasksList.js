import React, { Component } from "react";
import axios from "axios";

let endpoint = window.location.href;

export default class TasksList extends Component {
    constructor(props) {
        super(props);

        this.state = {
            data: [],
            err: null
        };
    }

    componentDidMount() {
        this.interval = setInterval(() => {
            this.GetTasks();
        }, 1000);
    }

    componentWillUnmount() {
        clearInterval(this.interval);
    }

    GetTasks = () => {
        axios.get(endpoint + "api/tasks").then((res) => {
            if (res.data) {
                this.setState({ data: res.data, err: null })
                //console.log(res.data)
            }
        }).catch(error => {
            this.setState({ data: [], err: error })
        });
    }

    render() {

        const headerStyle = {
            padding: "10px",
            backgroundColor: "#DADBDD",
            fontFamily: "monospace",
            fontSize: 15
        }

        const dataNormalStyle = {
            paddingRight: "10px",
            paddingLeft: "10px",
            fontFamily: "monospace",
            backgroundColor: "#D5F5E3",
            fontSize: 14
        }

        const dataOutdatedStyle = {
            paddingRight: "10px",
            paddingLeft: "10px",
            fontFamily: "monospace",
            backgroundColor: "#FF6347",
            fontSize: 14
        }

        const dataDoneStyle = {
            paddingRight: "10px",
            paddingLeft: "10px",
            fontFamily: "monospace",
            backgroundColor: "#50C878",
            fontSize: 14
        }

        const errorStyle = {
            fontFamily: "monospace",
            fontWeight: 1,
            color: "red",
            fontSize: 15
        }

        function taskStyle(task) {
            if (task.Solution !== null) { return dataDoneStyle }
            if (Date.parse(new Date()) > Date.parse(stringToTime(task.TimeoutAt))) { return dataOutdatedStyle }
            return dataNormalStyle
        }

        function datesDiffToString(bigger, lower) {
            var biggerNumber = Date.parse(bigger)
            var lowerNumber = Date.parse(lower)
            if (lowerNumber > biggerNumber) { return "00:00:00" }
            return (new Date(biggerNumber - lowerNumber)).toISOString().split('T')[1].split('.')[0]
        }

        function taskElapsed(task) {
            if (task.Solution !== null) {
                return datesDiffToString(task.Solution.FoundedAt, task.StartedAt)
            }
            return datesDiffToString(new Date(), task.StartedAt)
        }

        function taskConfirmed(task) {
            if (task.Solution !== null) { return "done" }
            if (task.ConfirmationsCounter > 0) {
                return datesDiffToString(new Date(), task.ConfirmedAt)
            }
            return ""
        }

        function taskTimeoutedOn(task) {
            if (task.Solution !== null) { return "done" }
            if (Date.parse(new Date()) > Date.parse(stringToTime(task.TimeoutAt))) {
                return datesDiffToString(new Date(), task.TimeoutAt)
            }
        }

        function stringToTime(timeField) {
            return new Date(Date.parse(timeField))
        }

        function axiosError2Text(error) {
            return JSON.stringify({
                code: error.code,
                data: error.response.data,
                message: error.message,
                status: error.response.status,
                statusText: error.response.statusText
            }, null, 4)
        }


        if (this.state.err !== null) {
            return (
                <div>
                    <p style={errorStyle}><pre>{axiosError2Text(this.state.err)}</pre></p>
                </div>
            )
        }


        return (
            <div>
                <h2>
                    Tasks List:
                </h2>
                <table border="1px" style={{ "borderCollapse": "collapse" }}>
                    <thead>
                        <tr>
                            <th style={headerStyle}>#</th>
                            <th style={headerStyle}>Number</th>
                            <th style={headerStyle}>Sequence</th>
                            <th style={headerStyle}>Deviation<br />threshold</th>

                            <th style={headerStyle}>Agent</th>
                            <th style={headerStyle}>Started At</th>
                            <th style={headerStyle}>Elapsed</th>
                            <th style={headerStyle}>Confirmed<br />Ago</th>
                            <th style={headerStyle}>Confirm<br />Count</th>
                            <th style={headerStyle}>Timeout<br />(Sec)</th>
                            <th style={headerStyle}>Timeouted<br />On</th>
                        </tr>
                    </thead>
                    <tbody>
                        {
                            this.state.data.map((task, index) => (
                                <tr key={index}>
                                    <td style={taskStyle(task)}>{index + 1}</td>
                                    <td style={taskStyle(task)}>{task.Number}</td>
                                    <td style={taskStyle(task)}>{task.Sequence}</td>
                                    <td style={taskStyle(task)}>{task.DeviationThreshold}</td>
                                    <td style={taskStyle(task)}>{task.Agent}</td>
                                    <td style={taskStyle(task)}>{(stringToTime(task.StartedAt)).toLocaleString()}</td>
                                    <td style={taskStyle(task)}>{taskElapsed(task)}</td>
                                    <td style={taskStyle(task)}>{taskConfirmed(task)}</td>
                                    <td style={taskStyle(task)}>{task.ConfirmationsCounter}</td>
                                    <td style={taskStyle(task)}>{task.TimeoutSec}</td>
                                    <td style={taskStyle(task)}>{taskTimeoutedOn(task)}</td>
                                </tr>
                            ))
                        }
                    </tbody>
                </table>
            </div>

        );
    }
}