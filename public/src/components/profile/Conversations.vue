<template>
  <v-container class="third-profile-step" fluild grid-list-md>
    <v-flex
      v-if="!local_conversations || local_conversations.length == 0"
      xs12
      text-sm-center
      text-xs-center
    >
      <v-alert
        color="info"
        flat
        value="You don't have any conversations"
      >You don't have any conversations</v-alert>
    </v-flex>
    <v-layout row wrap>
      <v-flex v-if="local_conversations">
        <v-list>
          <v-list-group
            v-for="(c) in local_conversations"
            :key="`conversation-${c.ID}`"
            no-action
            prepend-icon="message"
          >
            <v-divider :key="c.ID"></v-divider>
            <v-list-tile slot="activator">
              <!-- <v-list-tile-avatar>
                  <img v-if="c.From.Avatar" :src="c.From.Avatar" :alt="c.From.Avatar">
              </v-list-tile-avatar>-->
              <!-- <v-list-tile-avatar>
                <img v-if="c.To.Avatar" :src="c.To.Avatar" :alt="c.To.Avatar">
              </v-list-tile-avatar>-->
              <v-chip
                v-if="c.FromID == connectedUserID"
                small
                color
                class="subheading"
              >{{ c.To.Username || c.To.Firstname || c.From.Lastname }}</v-chip>
              <v-chip
                v-if="c.FromID != connectedUserID"
                small
                color
                class="subheading"
              >{{ c.From.Username || c.From.Firstname || c.From.Lastname }}</v-chip>
              <v-list-tile-sub-title
                class="text--primary"
              >Last message : {{ c.Messages[c.Messages.length - 1].Date| formatDate('MM/DD/YYYY') }} at {{ c.Messages[c.Messages.length - 1].Date| formatDate("HH:mm") }}</v-list-tile-sub-title>

              <!-- <v-list-tile-title>{{ c.To.Username || c.To.Firstname || c.To.Lastname }}</v-list-tile-title> -->
              <v-list-tile-action>
                <v-layout row>
                  <v-flex>
                    <v-btn color="primary" flat @click.prevent="deleteConversation(c)">
                      <v-icon>delete</v-icon>
                    </v-btn>
                  </v-flex>
                </v-layout>
              </v-list-tile-action>
            </v-list-tile>

            <v-list-tile v-for="m in c.Messages" :key="`message-${m.ID}`">
              <v-list-tile-avatar v-if="m.FromID == connectedUserID">
                <img v-if="c.To.Avatar" :src="c.To.Avatar" :alt="c.To.Avatar">
              </v-list-tile-avatar>
              <v-list-tile-avatar v-if="m.FromID != connectedUserID">
                <img v-if="c.To.Avatar" :src="c.From.Avatar" :alt="c.From.Avatar">
              </v-list-tile-avatar>
              <v-list-tile-content>
                <v-list-tile-sub-title
                  v-if="m.FromID != connectedUserID"
                  class
                >{{ c.From.Username }}:</v-list-tile-sub-title>
                <v-list-tile-sub-title v-if="m.FromID == connectedUserID" class>You:</v-list-tile-sub-title>
                <v-list-tile-title class="body-1">{{ m.Text }}</v-list-tile-title>
              </v-list-tile-content>

              <v-list-tile-action v-if="m.FromID != connectedUserID">
                <v-layout row>
                  <v-flex>
                    <v-btn color="primary" flat @click.prevent="openMessageDialog(c)">
                      <v-icon>reply</v-icon>
                    </v-btn>
                  </v-flex>
                </v-layout>
              </v-list-tile-action>
            </v-list-tile>
          </v-list-group>
        </v-list>
      </v-flex>
    </v-layout>
    <v-layout row justify-center>
      <v-dialog
        v-if="local_conversations"
        id="contact-dialog"
        v-model="showContactDialog"
        width="500"
      >
        <v-card>
          <v-toolbar color="primary">
            <v-card-title class="title font-weight-regular">Write your reply</v-card-title>
          </v-toolbar>
          <v-form v-model="messageFormValid">
            <v-card-text>
              <v-text-field
                v-if="!email"
                name="Email"
                label="Your email"
                autocomplete="email"
                v-model="message.Email"
                :rules="emailRules"
                autofocus
              ></v-text-field>
              <v-textarea
                name="Message"
                label="Your Message"
                v-model="message.Text"
                :rules="textRules"
                row="1"
                maxlength="128"
                hide-details
                no-resize
                autofocus
              ></v-textarea>
            </v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="primary" flat @click.prevent.native="showContactDialog = false">Cancel</v-btn>
              <v-btn color="primary" flat @click.native="reply" :disabled="!messageFormValid">Send</v-btn>
            </v-card-actions>
          </v-form>
        </v-card>
      </v-dialog>
    </v-layout>
    <app-snack-bar :state="snackbar" @snackClose="snackbar = false" :text="snackbarText"></app-snack-bar>
  </v-container>
</template>


<script>
import AppSnackBar from "@/components/utils/AppSnackBar";
import conversationRepo from "../../repositories/conversation.js";
import { mapState } from "vuex";

export default {
  name: "Conversations",
  props: ["conversations"],
  components: { AppSnackBar },
  data() {
    return {
      local_conversations: this.conversations.map(c => ({ ...c })),
      snackbar: false,
      snackbarTimeout: 3000,
      snackbarText: "your conversation has been successfully deleted",
      focusedConversation: null,
      showContactDialog: false,
      messageFormValid: false,
      emailRules: [
        v => !!v || "E-mail is required",
        v => /.+@.+/.test(v) || "E-mail must be valid"
      ],

      textRules: [
        v => !!v || "Message is required",
        v => (v && v.length >= 5) || "Message must be more than 20 characters"
      ]
    };
  },
  computed: {
    ...mapState({
      email: state => state.auth.email,
      connectedUserID: state => state.user.profile.ID
    }),
    message() {
      return { FromID: null, ToID: null, Email: this.email, Text: "" };
    }
  },
  watch: {
    snackbar(v) {
      if (!v) return;
      var that = this;
      setTimeout(function() {
        that.snackbar = false;
      }, that.snackbarTimeout);
    }
  },
  methods: {
    openMessageDialog: function(c) {
      this.showContactDialog = true;
      this.message.ToID = c.FromID == this.connectedUserID ? c.ToID : c.FromID;
      this.focusedConversation = c;
    },
    reply: function(e) {
      conversationRepo
        .sendMessage(this.message)
        .then(({ data }) => {
          this.snackbarText = "Your messages has been sent";
          this.snackbar = true;
          this.showContactDialog = false;
          data.Date = new Date();
          this.focusedConversation.Messages.push(data);
        })
        .catch(res => {
          this.snackbarText = "An error occured while sending your message";
          this.snackbar = true;
          this.showContactDialog = false;
        });
    },
    deleteConversation(c) {
      var that = this;

      this.focusedConversation = c;

      console.log(c, this.focusedConversation);
      if (c.ID != null) {
        conversationRepo
          .delete(c.ID)
          .then(function({ data }) {
            console.log(data, that.local_conversations);
            that.local_conversations = that.local_conversations.filter(function(
              c
            ) {
              console.log(that.focusedConversation.ID, c.ID);
              return that.focusedConversation.ID != c.ID;
            });
            that.snackbarText =
              "this conversation has been successfully deleted";
            that.snackbar = true;
          })
          .catch(() => {
            that.snackbarText = "there was an error deleting this conversation";
            that.snackbar = true;
          });
      }
    }
  }
};
</script>

<style lang="scss">
// .page-line:hover {
//   background: rgba($color: #607d8b, $alpha: 0.12);
// }
</style>

