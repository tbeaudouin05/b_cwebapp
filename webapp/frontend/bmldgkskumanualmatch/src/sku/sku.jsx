import React, { Component } from 'react';

export default class ShowBamiloSKU extends Component {
    constructor(props) {
        super(props);
        this.state = { value: '' };
        this.handleSubmit = this.handleSubmit.bind(this);
        this.handleChange = this.handleChange.bind(this);
    }

    handleChange(event) {
        this.setState({ value: event.target.value });
        event.preventDefault();
    }

    handleSubmit(event) {
        this.setState({ value: event.target.value });
        const { UpdateFrequency } = this.props;
        UpdateFrequency(this.input.current.value);
        event.preventDefault();
    }
    render() {
        return (

            <div className="row">
                <div className="img__wrap">
                    <img src={this.props.BmlImgLink} alt={this.props.BmlSkuName} className="img-thumbnail" height="150" width="150" />
                    <p className="img__description">{this.props.BmlSkuName}</p>
                </div>
                <form className="row_button" action={this.props.BmlSkuLink} target="_blank">
                    <button type="submit" className="btn btn-primary mb-2" >See on Bamilo</button>
                </form>
                <div>
                <p >
                <div className="row-gap-sm"></div>
                    Price: <strong  >{this.props.BmlPrice}</strong>
                </p>
                </div>




            </div>


        )
    }
}