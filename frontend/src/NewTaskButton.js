import React, { Component, useEffect } from "react";
import axios from "axios";
import { Card, Header, Form, Input, Icon } from "semantic-ui-react";

let endpoint = window.location.href;

class NewTaskButton extends Component {
    constructor(props) {
        super(props);
        this.state = {};
    }


    Create = () => {
        axios.get(endpoint + "api/task/new").then((res) => {
            if (res.data) {
                console.log(res.data)
            }
        })
    }

    render() {
        return (
            <div>
                <button onClick={this.Create}>
                    Create New Task
                </button>
            </div>
        );
    }
}

export default NewTaskButton;