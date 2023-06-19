import React, { Component } from "react";
import axios from "axios";

let endpoint = window.location.href;

export default class NewTaskButton extends Component {
    constructor(props) {
        super(props);
        this.state = {};
    }


    Create = () => {
        axios.get(endpoint + "api/task/console").then((res) => {
            if (res.data) {
                console.log(res.data)
            }
        }).catch(error => { });
    }

    render() {
        return (
            <div>
                <button onClick={this.Create}>
                    Create New Fake Task
                </button>
            </div>
        );
    }
}