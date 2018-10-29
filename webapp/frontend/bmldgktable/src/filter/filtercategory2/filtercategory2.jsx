import React, { Component } from 'react';
import OptionFilterCategory2 from '../filteroption/option.jsx'


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
            <form className="form-inline">
                <div className="form-group row">
                    <div className="col-sm-7">
                        <select if="selectCategory2" className="form-control" value={this.state.value} onChange={this.handleChange}>
                            {this.props.optionList.map(option => {
                                return (
                                    <OptionFilterCategory2 optionValue={option.optionValue} optionText={option.optionText} />
                                )
                            })}
                        </select>
                    </div>
                </div>

            </form>
        )
    }
}

