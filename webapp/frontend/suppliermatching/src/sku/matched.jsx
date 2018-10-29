import React, { Component } from 'react';

export default class ShowMatchedSKU extends Component {
    constructor(props) {
        super(props);

    }
    onClick(e) {
        e.preventDefault();
        const { ApplyUnmatch } = this.props;
        ApplyUnmatch();
    }

    render() {
        return (
            <div>
                <div>
                    <strong>Matched Product:</strong>
                    <div className="img__wrap">
                        <img src={this.props.DgkImgLink} height="170" width="170" />
                    </div>
                    <li>
                </li>
                <li>
                </li>
                <li>
                </li>
                    <form className="row_button" onClick={this.onClick.bind(this)} target="_blank">
                        <button type="submit" className="btn btn-primary mb-2" >Unmatch</button>
                    </form>
                </div>
                
            </div>
        )
    }
}
