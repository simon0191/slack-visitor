import Vue from 'vue';
import { Component, watch } from 'vue-property-decorator';
import { Link } from './link';
import { Logger } from '../../util/log';

@Component({
    template: require('./navbar.html')
})
export class NavbarComponent extends Vue {

    public inverted: boolean = true;
    protected logger: Logger;

    links: Link[] = [
        new Link('Github', 'https://github.com/simon0191/slack-visitor'),
    ];

    @watch('$route.path')
    pathChanged() {
        this.logger.info('Changed current path to: ' + this.$route.path);
    }

    mounted() {
        if (!this.logger) this.logger = new Logger();
    }
}
