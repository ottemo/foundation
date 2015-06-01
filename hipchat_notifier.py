import httplib, urllib
import getopt, sys, os
import subprocess

def get_connection():
      return httplib.HTTPSConnection('api.hipchat.com')

def get_url(token):
  return '/v1/rooms/message?auth_token=%s' % token

def get_data_from_git(format_string, commit):
    return subprocess.check_output(['git', 'log', '-1', '--format=format:%s' % format_string, commit])

def get_author(commit):
    return get_data_from_git('%an <%ae>', commit)

def get_date(commit):
    return get_data_from_git('%aD', commit)

def get_title(commit):
    return get_data_from_git('%s', commit)

def get_full_message(commit):
    return get_data_from_git('%b', commit)

def post_message(connection, url, room, success, project):
    headers = {'Content-type': 'application/x-www-form-urlencoded'}
    build_url = os.environ['BUILD_URL']
    build_number = os.environ['BUILD_NUMBER']
    branch = os.environ['BRANCH']
    commit = os.environ['COMMIT']

    status_text = 'succeeded' if success else 'failed'
    color = 'green' if success else 'red'
    if branch != 'develop':
        title = '<a href="%s">Build #%s</a> %s for project <strong>%s</strong> on branch %s' % (build_url, build_number, status_text, project, branch)
    elif branch == 'develop':
        title = '<a href="%s">Build & Deployment #%s</a> %s for project <strong>%s</strong> on branch %s' % (build_url, build_number, status_text, project, branch)
    author = '<strong>%s</strong><br>%s' % (get_author(commit), get_date(commit))
    description = '<strong>%s</strong><br>%s' % (get_title(commit), get_full_message(commit))
    message = '<br>'.join([title, author, description])

    message = {
      'room_id': room,
      'from': 'Shippable',
      'color': color,
      'message': message,
      'notify': True,
    }

    connection.request('POST', url, urllib.urlencode(message), headers)
    response = connection.getresponse()
    print response.read().decode()

def main():
    try:
        opts, args = getopt.getopt(sys.argv[1:], ':sf', ['project=', 'room=', 'token='])
    except getopt.GetoptError as err:
        print str(err)
        sys.exit(2)

    success = False
    room = None
    project = None
    room = None
    token = None
    for o, arg in opts:
        if o == '-s':
            success = True
        elif o == '--project':
            project = arg
        elif o == '--room':
            room = arg
        elif o == '--token':
            token = arg

    connection = get_connection()
    url = get_url(token)
    post_message(connection, url, room, success, project)

if __name__ == '__main__':
    main()
