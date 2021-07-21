import { h, Component } from '../../../deps/preact';
import { DetailView } from './detailview';
import { IconView } from './iconview';
import { ListView } from './listview';
import { File } from './model';
import { TileView } from './tileview';

export const enum ViewMode {
    Icons,
    List,
    Details,
    Tiles,
}

interface Props {
    files: File[];
    mode: ViewMode;
};

export class Explorer extends Component<Props> {
    render() {
        switch(this.props.mode) {
        case ViewMode.Icons:
            return <IconView files={this.props.files} />;
        case ViewMode.List:
            return <ListView files={this.props.files} />;
        case ViewMode.Details:
            return <DetailView files={this.props.files} />;
        case ViewMode.Tiles:
            return <TileView files={this.props.files} />;
        }
    }
};
