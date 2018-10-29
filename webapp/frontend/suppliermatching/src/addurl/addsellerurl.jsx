import React, { Component } from 'react';


export default class SellerUrl extends Component {

    // cf. this: https://reactjs.org/docs/forms.html

    constructor(props) {
        super(props);
        this.state = { value: "" };

        this.handleChange = this.handleChange.bind(this);
    }

    handleChange(event) {
        this.setState({ value: event.target.value });
       
        event.preventDefault();
    }

    onClick(e) {
        e.preventDefault();
        const { ApplySellerUrl } = this.props;
        ApplySellerUrl(this.state.value);
    }

    render() {
        return (
            <div className="row">
                <form>
                    <p > </p>
                    <p ><strong> Add new Digikala supplier URL: </strong> </p>
                    <div className="form-group mx-sm-3 mb-2">
                        <label htmlFor="search" className="sr-only"></label>
                        <input type="text" className="form-control" id="url" placeholder="https://www.digikala.com/seller/00/" value={this.state.value} onChange={this.handleChange} />
                    </div>
                    <form className="row_button" onClick={this.onClick.bind(this)} target="_blank">
                        <button type="submit" className="btn btn-primary mb-2">Submit</button>
                    </form>
                    <form className="row" ValidationStatus={this.props.ValidationStatus} >
                    <td key='4'> {this.props.ValidationStatus} </td>

                    </form>
                </form>
            </div>
        )
    }
}

