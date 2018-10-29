import React, { Component } from 'react';
import OptionFilterCategory2 from './option2'


export default class FilterCategory2 extends Component {

    // cf. this: https://reactjs.org/docs/forms.html

    constructor(props) {
        super(props);
        this.state = { value: '' };

        this.handleChange = this.handleChange.bind(this);
    }

    handleChange(event) {
        this.setState({ value: event.target.value });
        this.props.chooseCategory2(event.target.value);
        event.preventDefault();
    }


    render() {
        return (
            <form>
                <select className="form-control" value={this.state.value} onChange={this.handleChange}>
                    {this.props.optionList2.map(option => {
                        return (
                            <OptionFilterCategory2 optionValue={option.optionValue} optionText={option.optionText} />
                        )
                    })}
                </select>
            </form>
        )
    }
}

