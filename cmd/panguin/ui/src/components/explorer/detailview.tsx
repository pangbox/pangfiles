import { h, Component } from '../../../deps/preact';
import { File } from './model';

interface Props {
    files: File[];
};

interface State {
    colWidths: number[];
    selectedIndex: number | null;
}

interface HeaderProps {
    colWidths: number[];
    setColWidth: (index: number, value: number) => void;
};

interface RowProps {
    file: File;
    colWidths: number[];
    selected: boolean;
    onClick: () => void;
};

interface DividerProps {
    width: number;
    setWidth: (value: number) => void;
}

export class Divider extends Component<DividerProps> {
    onMouseDown = (event: MouseEvent) => {
        document.addEventListener("mousemove", this.onMouseMove);
        document.addEventListener("mouseup", this.onMouseUp);
    };

    onMouseMove = (event: MouseEvent) => {
        this.props.setWidth(this.props.width + event.movementX / window.devicePixelRatio);        
    };

    onMouseUp = (event: MouseEvent) => {
        document.removeEventListener("mousemove", this.onMouseMove);
        document.removeEventListener("mouseup", this.onMouseUp);
    };

    componentWillUnmount() {
        document.removeEventListener("mousemove", this.onMouseMove);
        document.removeEventListener("mouseup", this.onMouseUp);
    }

    render() {
        return (
            <div class="explorer-detail-divider" onMouseDown={this.onMouseDown}>
                <div class="explorer-detail-divider-line"></div>
            </div>
        );
    }
}

export class DetailHeader extends Component<HeaderProps> {
    render() {
        return <div class="explorer-detail-row explorer-detail-header">
            <div class="explorer-detail-cell" style={{width: `${this.props.colWidths[0]}px`}}>
                Name
                <Divider width={this.props.colWidths[0]} setWidth={w => this.props.setColWidth(0, w)} />
            </div>
            <div class="explorer-detail-cell" style={{width: `${this.props.colWidths[1]}px`}}>
                Size
                <Divider width={this.props.colWidths[1]} setWidth={w => this.props.setColWidth(1, w)} />
            </div>
        </div>;
    }
}

export class DetailRow extends Component<RowProps> {
    onClick = (event: MouseEvent) => {
        event.stopPropagation();
        this.props.onClick();
    }
    render() {
        return <div class={`explorer-detail-row ${this.props.selected ? "selected": ""}`}>
            <div class="explorer-detail-cell" style={{width: `${this.props.colWidths[0]}px`}} onClick={this.onClick}>
                {this.props.file.name}
            </div>
            <div class="explorer-detail-cell" style={{width: `${this.props.colWidths[1]}px`}} onClick={this.onClick}>
                {this.props.file.size}
            </div>
        </div>;
    }
}

export class DetailView extends Component<Props, State> {
    constructor(props?: Props, context?: any) {
        super(props, context);

        this.state = {
            colWidths: [200, 50],
            selectedIndex: null,
        };
    }

    render() {
        return <div class="explorer-detail-view">
            <DetailHeader
                colWidths={this.state.colWidths}
                setColWidth={(idx, w) => {
                    const colWidths = [...this.state.colWidths];
                    colWidths[idx] = Math.max(w, 50);
                    this.setState({colWidths});
                }}/>
            <div
                class="explorer-detail-pane"
                onClick={() => this.setState({ selectedIndex: null })}>
                {this.props.files.map((file, i) => (
                    <DetailRow
                        file={file}
                        colWidths={this.state.colWidths}
                        key={file.name}
                        selected={this.state.selectedIndex === i}
                        onClick={() => this.setState({ selectedIndex: i })} />
                ))}
            </div>
        </div>;
    }
};
