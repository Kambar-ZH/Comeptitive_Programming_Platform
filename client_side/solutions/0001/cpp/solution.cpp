#include <bits/stdc++.h>

using namespace std;

#define ll                  long long
#define vt                  vector
#define For(i, n)           for(int i = 0; i < n; ++i)
#define pii                 pair<int, int>

const int MOD = 1e9 + 7;

int main()
{
    int n; cin >> n;
    int sum = 0;
    For(i, n) {
        int x; cin >> x;
        sum += x;
    }
    cout << sum << endl;
    return 0;
}