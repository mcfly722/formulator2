import React, { Component } from "react";
import BestSolutionsList from "./BestSolutionsList.js";
import axios from "axios";

let endpoint = window.location.href;

export default class State extends Component {

    constructor(props) {
        super(props);

        this.state = {
            state: null,
            err: null
        };
    }

    componentDidMount() {
        this.interval = setInterval(() => {
            this.GetState();
        }, 1000);
    }

    GetState() {
        axios.get(endpoint + "api/state").then((res) => {
            if (res.data) {
                this.setState({ state: res.data, err: null })
                //this.bestSolutionList.setState({  })
                //console.log(res.data)
            }
        }).catch(error => {
            this.setState({ state: {}, err: error })

        });
    }


    render() {

        const errorStyle = {
            fontFamily: "monospace",
            fontWeight: 1,
            color: "red",
            fontSize: 15
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


        if (typeof this.state.state === "undefined" || this.state.state === null) {
            return (<div></div>)
        }

        return (
            <div>
                <h2>
                    State:
                    <BestSolutionsList solutions={this.state.state.Solutions} />
                </h2>
            </div>
        )
    };

}
