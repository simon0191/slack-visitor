import Vue from 'vue';
import { Component, watch } from 'vue-property-decorator';
import { Logger } from '../../util/log';

@Component({
    template: require('./navbar.html')
})
export class NavbarComponent extends Vue {

    public inverted: boolean = true;
    protected logger: Logger;

    @watch('$route.path')
    pathChanged() {
        this.logger.info('Changed current path to: ' + this.$route.path);
    }

    mounted() {
        if (!this.logger) this.logger = new Logger();
    }
}
