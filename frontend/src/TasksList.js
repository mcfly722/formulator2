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
            if (task.Done === true) { return dataDoneStyle }
            if (task.TimeoutedOnSec > 0) { return dataOutdatedStyle }
            return dataNormalStyle
        }

        function taskDoneValue(task, done, value) {
            if (task.Done) {
                return done
            }
            return value
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
                    <thead >
                        <tr >
                            <th style={headerStyle}>Number</th>
                            <th style={headerStyle}>Sequence</th>
                            <th style={headerStyle}>Agent</th>
                            <th style={headerStyle}>Started At</th>
                            <th style={headerStyle}>Elapsed</th>
                            <th style={headerStyle}>Confirmation</th>
                            <th style={headerStyle}>Timeouted On (Sec)</th>
                        </tr>
                    </thead>
                    <tbody>
                        {
                            this.state.data.map((task, index) => (
                                <tr key={index}>
                                    <td style={taskStyle(task)}>{task.Number}</td>
                                    <td style={taskStyle(task)}>{task.Sequence}</td>
                                    <td style={taskStyle(task)}>{task.Agent}</td>
                                    <td style={taskStyle(task)}>{task.StartedAt}</td>
                                    <td style={taskStyle(task)}>{task.Elapsed}</td>
                                    <td style={taskStyle(task)}>{taskDoneValue(task, "done", task.LastConfirmationAgo)}</td>
                                    <td style={taskStyle(task)}>{taskDoneValue(task, "done", task.TimeoutedOnSec)}</td>
                                </tr>
                            ))
                        }
                    </tbody>
                </table>
            </div>

        );
    }
}